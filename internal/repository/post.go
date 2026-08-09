package repository
import (
	"context"
	"encoding/json"
	"fmt"

	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
)

type postRepo struct {
	db DBTX
}

func NewPostRepo(db DBTX) *postRepo {
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

const getPostsQuery = `
SELECT P.id, P.title, P.content, likes_count, dislikes_count, P.created_at, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at, COALESCE(R.id, '00000000-0000-0000-0000-000000000000'), COALESCE(R.is_like, FALSE),` + postCategoriesColumn + `
FROM posts P
JOIN users U
ON P.author_id = U.id
LEFT JOIN reactions R
ON P.id = R.parent_id AND R.author_id = ?
ORDER BY P.created_at DESC
LIMIT ? OFFSET ?
`

// categoryRow mirrors the json_object keys above. Decoding through it keeps the
// SQL free of any knowledge of how crypto.UUID is laid out in Go.
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
		id, err := crypto.ParseUUID(row.ID)
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
		&post.ID.Value,
		&post.Title,
		&post.Content,
		&post.Likes,
		&post.DisLikes,
		&post.CreatedAt,
		&post.Author.ID.Value,
		&post.Author.NickName,
		&post.Author.LastName,
		&post.Author.FirstName,
		&post.Author.Gender,
		&post.Author.Age,
		&post.Author.CreatedAt,
		&post.LikeInfo.ID.Value,
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

	if _, err := tx.ExecContext(ctx, savePostQuery, post.ID.Value, post.AuthorID.Value, post.Title, post.Content); err != nil {
		return sqlite.TranslateError(err)
	}

	for _, categoryID := range post.Categories {
		if _, err := tx.ExecContext(ctx, savePostCategoryQuery, post.ID.Value, categoryID.Value); err != nil {
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

func (p *postRepo) GetPost(ctx context.Context, userID, postID crypto.UUID) (*domain.PostInfo, error) {
	row := p.db.QueryRowContext(ctx, getPostQuery, userID.Value, postID.Value)
	return scanPost(row)
}


//const getPostsQuery = `SELECT id, author_id, title, content, created_at, updated_at FROM posts 
//						ORDER BY created_at DESC 
//						LIMIT ? OFFSET ?; `

func (p *postRepo) GetPosts(ctx context.Context, userID crypto.UUID, limit, offset int) ([]*domain.PostInfo, error) {
	rows, err := p.db.QueryContext(ctx, getPostsQuery, userID.Value, limit, offset)	
	
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
	return posts, nil 
}

const deletePostQuery = `DELETE FROM posts WHERE id = ?`

func (p *postRepo) DeletePost(ctx context.Context, postID crypto.UUID) error {
	_, err := p.db.ExecContext(ctx, deletePostQuery, postID.Value)
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepo) ExecPatchQuery(ctx context.Context, query string, postID crypto.UUID, args []any) error {
	result, err := p.db.ExecContext(ctx, query, append(args, postID.Value)...) 	
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected() 
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected edited rows")
	}
	return nil	
}

