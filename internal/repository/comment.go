package repository

import (
	"context"
	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
	"uuid"
)

type commentRepo struct {
	db DBTX
}

func NewCommentRepo(db DBTX) *commentRepo {
	return &commentRepo{db: db}
}

const getCommentsQuery = `
SELECT C.id,  C.content, C.created_at, C.post_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at
FROM comments C
JOIN users U
ON C.author_id = U.id
WHERE C.post_id = ?
`

func scanComment(row scanner) (*domain.CommentInfo, error) {
	comment := domain.CommentInfo{}
	err := row.Scan(
		&comment.ID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.PostID,
		&comment.Author.ID,
		&comment.Author.NickName,
		&comment.Author.FirstName,
		&comment.Author.LastName,
		&comment.Author.Gender,
		&comment.Author.Age,
		&comment.Author.CreatedAt)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return &comment, nil
}

const saveCommentQuery = `INSERT INTO comments (id, post_id, author_id, content) VALUES(?, ?, ?, ?)`

func (p *commentRepo) SaveComment(ctx context.Context, comment domain.Comment) error {
	_, err := p.db.ExecContext(ctx, saveCommentQuery, comment.ID.String(), comment.PostID.String(), comment.AuthorID.String(), comment.Content)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const getCommentQuery = `
SELECT C.id,  C.content, C.created_at, C.post_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at
FROM comments C
JOIN users U
ON C.author_id = U.id
WHERE C.id = ?
`

func (p *commentRepo) GetComment(ctx context.Context, userID, commentID uuid.UUID) (*domain.CommentInfo, error) {
	row := p.db.QueryRowContext(ctx, getCommentQuery, commentID.String())
	return scanComment(row)
}

func (p *commentRepo) GetComments(ctx context.Context, userID, postID uuid.UUID) ([]*domain.CommentInfo, error) {
	rows, err := p.db.QueryContext(ctx, getCommentsQuery, postID.String())

	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	defer rows.Close()
	var comments []*domain.CommentInfo
	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return comments, nil
}


const deleteCommentQuery = `DELETE FROM comments WHERE id = ?`

func (p *commentRepo) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	_, err := p.db.ExecContext(ctx, deleteCommentQuery, commentID.String())
	if err != nil {
		return err
	}
	return nil
}
