package repository


import (
	"fmt"
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
)

type commentRepo struct {
	db DBTX
}

func NewCommentRepo(db DBTX) *commentRepo {
	return &commentRepo{db: db}
}


const getCommentsQuery = `
SELECT C.id,  C.content, C.created_at, C.parent_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at, COALESCE(R.id, '00000000-0000-0000-0000-000000000000'), COALESCE(R.is_like, FALSE) 
FROM comments C  
JOIN users U 
ON C.author_id = U.id
LEFT JOIN  reactions R  
ON R.author_id = ? AND C.id = R.parent_id  
WHERE C.parent_id = ?
`

func scanComment(row scanner) (*domain.CommentInfo, error) {
	comment := domain.CommentInfo{}
	err := row.Scan(
		&comment.ID.Value,
		&comment.Content,
		&comment.CreatedAt,
		&comment.ParentID.Value,
		&comment.Author.ID.Value,
		&comment.Author.NickName,
		&comment.Author.FirstName,
		&comment.Author.LastName,
		&comment.Author.Gender,
		&comment.Author.Age,
		&comment.Author.CreatedAt,
		&comment.LikeInfo.ID.Value,
		&comment.LikeInfo.IsLike)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return &comment, nil
}

const saveCommentQuery = `INSERT INTO comments (id, parent_id, author_id, content) VALUES(?, ?, ?, ?)`

func (p *commentRepo) SaveComment(ctx context.Context, comment domain.Comment) error {
	_, err := p.db.ExecContext(ctx, saveCommentQuery, comment.ID.Value, comment.ParentID.Value, comment.AuthorID.Value, comment.Content)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const getCommentQuery = `
SELECT C.id,  C.content, C.created_at, C.parent_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at, COALESCE(R.id, '00000000-0000-0000-0000-000000000000'), COALESCE(R.is_like, FALSE)
FROM comments C  
JOIN users U 
ON C.author_id = U.id
LEFT JOIN  reactions R  
ON R.author_id = ? AND C.id = R.parent_id  
WHERE C.id = ?
`


func (p *commentRepo) GetComment(ctx context.Context, userID, commentID crypto.UUID) (*domain.CommentInfo, error) {
	row := p.db.QueryRowContext(ctx, getCommentQuery, userID.Value, commentID.Value)	
	return scanComment(row) 
}


//const getCommentsQuery = `SELECT id, author_id, title, content, created_at, updated_at FROM comments 
//						ORDER BY created_at DESC 
//						LIMIT ? OFFSET ?; `

func (p *commentRepo) GetComments(ctx context.Context, userID, parentID crypto.UUID) ([]*domain.CommentInfo, error) {
	rows, err := p.db.QueryContext(ctx, getCommentsQuery, userID.Value, parentID.Value)	
	
	if err != nil {
		return nil, sqlite.TranslateError(err)	
	}
	defer rows.Close()
	var comments []*domain.CommentInfo
	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			return nil, sqlite.TranslateError(err)	
		}
		comments = append(comments, comment)
	}
	return comments, nil 
}

const deleteCommentQuery = `DELETE FROM comments WHERE id = ?`

func (p *commentRepo) DeleteComment(ctx context.Context, commentID crypto.UUID) error {
	_, err := p.db.ExecContext(ctx, deleteCommentQuery, commentID.Value)
	if err != nil {
		return err
	}
	return nil
}

func (p *commentRepo) ExecPatchQuery(ctx context.Context, query string, commentID crypto.UUID, args []any) error {
	result, err := p.db.ExecContext(ctx, query, append(args, commentID.Value)...) 	
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

