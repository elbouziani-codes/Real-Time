package repository
import (
	"fmt"
	"context"
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


const getPostsQuery = `
SELECT P.id, P.title, P.content, likes_count, dislikes_count, P.created_at, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at, COALESCE(R.id, '00000000-0000-0000-0000-000000000000'), COALESCE(R.is_like, FALSE)
FROM posts P  
JOIN users U 
ON P.author_id = U.id
LEFT JOIN reactions R
ON P.id = R.parent_id AND R.author_id = ? 
ORDER BY P.created_at DESC 
LIMIT ? OFFSET ?
`

func scanPost(row scanner) (*domain.PostInfo, error) {
	post := domain.PostInfo{}
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
		&post.LikeInfo.IsLike)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return &post, nil
}

const savePostQuery = `INSERT INTO posts (id, author_id, title, content) VALUES(?, ?, ?, ?)`

func (p *postRepo) SavePost(ctx context.Context, post domain.Post) error {
	_, err := p.db.ExecContext(ctx, savePostQuery, post.ID.Value, post.AuthorID.Value, post.Title, post.Content)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const getPostQuery = `
SELECT P.id, P.title, P.content, likes_count, dislikes_count, P.created_at, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at, COALESCE(R.id, '00000000-0000-0000-0000-000000000000'), COALESCE(R.is_like, FALSE) 
FROM posts P  
JOIN users U 
ON P.author_id = U.id
LEFT JOIN reactions R
ON P.id = R.parent_id AND R.author_id = ?  
WHERE P.id = ?
`


func (p *postRepo) GetPost(ctx context.Context, userID, postID crypto.UUID) (*domain.PostInfo, error) {
	fmt.Println(userID, postID)
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

