package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
	"uuid"
)

type postRepo struct {
	db DB
}

func NewPostRepo(db DB) *postRepo {
	return &postRepo{db: db}
}

// postCategoriesColumn aggregates a post's categories into one JSON array. It
// is a correlated subquery rather than a join so it cannot multiply post rows
// against the reactions join below.
const postCategoriesColumn = `
COALESCE((
	SELECT json_group_array(json_object('id', C.id, 'title', C.title, 'icon', C.icon, 'created_at', C.created_at))
	FROM post_categories PC
	JOIN categories C ON PC.category_id = C.id
	WHERE PC.post_id = P.id
), '[]')`

// getPostsQueryHead and getPostsQueryTail sandwich the optional filter clauses
// built by buildPostFilter. Keeping them apart lets the filters be composed
// without string-formatting any caller-supplied value into the SQL.
const getPostsQueryHead = `
SELECT P.id, P.title, P.content, likes_count, dislikes_count, P.created_at, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at, COALESCE(R.id, '00000000-0000-0000-0000-000000000000'), COALESCE(R.is_like, FALSE),` + postCategoriesColumn + `
FROM posts P
JOIN users U
ON P.author_id = U.id
LEFT JOIN reactions R
ON P.id = R.parent_id AND R.author_id = ?
`

// getPostsQueryTail closes the listing. The id breaks ties on created_at, which
// only has one-second resolution: without a total order, two posts written in
// the same second could come back in either order on different pages, and the
// cursor would then skip or repeat them.
const getPostsQueryTail = `
ORDER BY P.created_at DESC, P.id DESC
LIMIT ?
`

const getPostCursorQuery = `SELECT created_at FROM posts WHERE id = ?`

// postCursorKey is the sort position of the post a client last received. The id
// belongs to it because created_at alone does not identify a single row.
type postCursorKey struct {
	createdAt int
	id        uuid.UUID
}

// resolvePostCursor reads the sort position of the cursor post. A cursor naming
// a post that has since been deleted is reported rather than silently returning
// an empty page, which a client could not tell apart from the end of the feed.
func (p *postRepo) resolvePostCursor(ctx context.Context, cursor uuid.UUID) (*postCursorKey, error) {
	if cursor == uuid.Nil() {
		return nil, nil
	}

	key := postCursorKey{id: cursor}
	err := p.db.QueryRowContext(ctx, getPostCursorQuery, cursor.String()).Scan(&key.createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.Error{Message: "unknown cursor", Code: domain.NotFoundCode}
	}
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return &key, nil
}

// buildPostFilter renders the filter and the cursor as SQL conditions plus their
// arguments. The category clause uses EXISTS rather than a join so a post in
// several of the requested categories is still returned once.
func buildPostFilter(filter domain.PostFilter, cursor *postCursorKey) (string, []any) {
	conditions := make([]string, 0, 3)
	args := make([]any, 0, len(filter.Categories)+3)

	if len(filter.Categories) > 0 {
		placeholders := strings.Repeat("?,", len(filter.Categories)-1) + "?"
		conditions = append(conditions, `EXISTS (
	SELECT 1 FROM post_categories PC
	WHERE PC.post_id = P.id AND PC.category_id IN (`+placeholders+`)
)`)
		for _, categoryID := range filter.Categories {
			args = append(args, categoryID.String())
		}
	}

	// R is already joined on this user, but it also matches dislikes; the
	// is_like test is what makes this "liked" rather than "reacted to".
	if filter.LikedOnly {
		conditions = append(conditions, "R.id IS NOT NULL AND R.is_like = TRUE")
	}

	// Keyset paging: keep the rows that fall strictly after the cursor under the
	// ORDER BY above. This condition comes last because its placeholders have to
	// follow the category ones in the argument list.
	if cursor != nil {
		conditions = append(conditions, "(P.created_at < ? OR (P.created_at = ? AND P.id < ?))")
		args = append(args, cursor.createdAt, cursor.createdAt, cursor.id.String())
	}

	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

// categoryRow mirrors the json_object keys above. Decoding through it keeps the
// SQL free of any knowledge of how uuid.UUID is laid out in Go.
type categoryRow struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	CreatedAt int    `json:"created_at"`
}

func decodeCategories(raw string) ([]domain.Category, error) {
	var rows []categoryRow
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		return nil, err
	}
	categories := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, domain.Category{
			ID:        id,
			Title:     row.Title,
			Icon:      row.Icon,
			CreatedAt: row.CreatedAt,
		})
	}
	return categories, nil
}

func scanPost(row scanner) (*domain.PostInfo, error) {
	post := domain.PostInfo{}
	var rawCategories string
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Likes,
		&post.DisLikes,
		&post.CreatedAt,
		&post.Author.ID,
		&post.Author.NickName,
		&post.Author.FirstName,
		&post.Author.LastName,
		&post.Author.Gender,
		&post.Author.Age,
		&post.Author.CreatedAt,
		&post.LikeInfo.ID,
		&post.LikeInfo.IsLike,
		&rawCategories)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	post.Categories, err = decodeCategories(rawCategories)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

const savePostQuery = `INSERT INTO posts (id, author_id, title, content) VALUES(?, ?, ?, ?)`

const savePostCategoryQuery = `INSERT INTO post_categories (post_id, category_id) VALUES(?, ?)`

// SavePost writes the post and its categories in one transaction. An unknown
// category id trips the post_categories foreign key, which rolls the whole
// insert back instead of leaving a post with no categories behind.
func (p *postRepo) SavePost(ctx context.Context, post domain.Post) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, savePostQuery, post.ID.String(), post.AuthorID.String(), post.Title, post.Content); err != nil {
		return sqlite.TranslateError(err)
	}

	for _, categoryID := range post.Categories {
		if _, err := tx.ExecContext(ctx, savePostCategoryQuery, post.ID.String(), categoryID.String()); err != nil {
			return sqlite.TranslateError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const getPostQuery = `
SELECT P.id, P.title, P.content, likes_count, dislikes_count, P.created_at, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at, COALESCE(R.id, '00000000-0000-0000-0000-000000000000'), COALESCE(R.is_like, FALSE),` + postCategoriesColumn + `
FROM posts P
JOIN users U
ON P.author_id = U.id
LEFT JOIN reactions R
ON P.id = R.parent_id AND R.author_id = ?
WHERE P.id = ?
`

func (p *postRepo) GetPost(ctx context.Context, userID, postID uuid.UUID) (*domain.PostInfo, error) {
	row := p.db.QueryRowContext(ctx, getPostQuery, userID.String(), postID.String())
	return scanPost(row)
}

func (p *postRepo) GetPosts(ctx context.Context, userID uuid.UUID, filter domain.PostFilter, limit int, cursor uuid.UUID) ([]*domain.PostInfo, error) {
	cursorKey, err := p.resolvePostCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	// Argument order has to track the query text: the reactions join owns the
	// first placeholder, then the filter's and the cursor's, then the limit.
	where, filterArgs := buildPostFilter(filter, cursorKey)
	query := getPostsQueryHead + where + getPostsQueryTail

	args := make([]any, 0, len(filterArgs)+2)
	args = append(args, userID.String())
	args = append(args, filterArgs...)
	args = append(args, limit)

	rows, err := p.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	defer rows.Close()
	var posts []*domain.PostInfo
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, sqlite.TranslateError(err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return posts, nil
}
