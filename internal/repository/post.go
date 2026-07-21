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

func scanPost(row scanner) (*domain.Post, error) {
	post := domain.Post{}
	err := row.Scan(
		&post.ID.Value,
		&post.AuthorID.Value,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt)
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

const getPostQuery = `SELECT id, author_id, title, content, created_at, updated_at FROM posts WHERE id = ?`

func (p *postRepo) GetPost(ctx context.Context, postID crypto.UUID) (*domain.Post, error) {
	row := p.db.QueryRowContext(ctx, getPostQuery, postID.Value)	
	return scanPost(row) 
}


const getPostsQuery = `SELECT id, author_id, title, content, created_at, updated_at FROM posts 
						ORDER BY created_at DESC 
						LIMIT ? OFFSET ?; `

func (p *postRepo) GetPosts(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	rows, err := p.db.QueryContext(ctx, getPostsQuery, limit, offset)	
	
	if err != nil {
		return nil, sqlite.TranslateError(err)	
	}
	defer rows.Close()
	var posts []*domain.Post
	for rows.Next() {
		fmt.Println("a")
		post, err := scanPost(rows)
		if err != nil {
			return nil, sqlite.TranslateError(err)	
		}
		fmt.Println(posts)
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

