package repository

import (
	"fmt"
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
)

type reactionRepo struct {
	db DBTX
}

func NewReactionRepo(db DBTX) *reactionRepo {
	return &reactionRepo{db: db}
}



func scanReaction(row scanner) (*domain.ReactionInfo, error) {
	reaction := domain.ReactionInfo{}
	err := row.Scan(
		&reaction.ID.Value,
		&reaction.IsLike,
		&reaction.CreatedAt,
		&reaction.ParentID.Value,
		&reaction.Author.ID.Value,
		&reaction.Author.NickName,
		&reaction.Author.FirstName,
		&reaction.Author.LastName,
		&reaction.Author.Gender,
		&reaction.Author.Age,
		&reaction.Author.CreatedAt)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return &reaction, nil
}

const saveReactionQuery = `INSERT INTO reactions (id, parent_id, author_id, is_like) VALUES(?, ?, ?, ?)`

func (p *reactionRepo) SaveReaction(ctx context.Context, reaction domain.Reaction) error {
	_, err := p.db.ExecContext(ctx, saveReactionQuery, reaction.ID.Value, reaction.ParentID.Value, reaction.AuthorID.Value, reaction.IsLike)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const getReactionQuery = `
SELECT R.id,  R.is_like, R.created_at, R.parent_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at
FROM reactions R  
JOIN users U 
ON R.author_id = U.id
WHERE R.id = ?
`


func (p *reactionRepo) GetReaction(ctx context.Context, reactionID crypto.UUID) (*domain.ReactionInfo, error) {
	row := p.db.QueryRowContext(ctx, getReactionQuery, reactionID.Value)	
	return scanReaction(row) 
}

const getReactionWithUserAndParentQuery = `
SELECT R.id,  R.is_like, R.created_at, R.parent_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at
FROM reactions R  
JOIN users U 
ON R.author_id = U.id
WHERE R.parent_id = ? AND R.author_id = ?
`
func (p *reactionRepo) GetReactionByUserAndParent(ctx context.Context, parentID, userID crypto.UUID) (*domain.ReactionInfo, error) {
	row := p.db.QueryRowContext(ctx, getReactionWithUserAndParentQuery, parentID.Value, userID.Value)	
	return scanReaction(row) 
}


const getReactionsQuery = `SELECT R.id,  R.is_like, R.created_at, R.parent_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at
						FROM reactions R
						JOIN users U
						ON R.author_id = U.id
						WHERE R.parent_id = ?
						ORDER BY R.created_at DESC`

func (p *reactionRepo) GetReactions(ctx context.Context, parentID crypto.UUID) ([]*domain.ReactionInfo, error) {
	rows, err := p.db.QueryContext(ctx, getReactionsQuery, parentID.Value)

	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	defer rows.Close()
	var reactions []*domain.ReactionInfo
	for rows.Next() {
		reaction, err := scanReaction(rows)
		if err != nil {
			return nil, sqlite.TranslateError(err)
		}
		reactions = append(reactions, reaction)
	}
	if err := rows.Err(); err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return reactions, nil
}

const deleteReactionQuery = `DELETE FROM reactions WHERE id = ?`

func (p *reactionRepo) DeleteReaction(ctx context.Context, reactionID crypto.UUID) error {
	_, err := p.db.ExecContext(ctx, deleteReactionQuery, reactionID.Value)
	if err != nil {
		return err
	}
	return nil
}
const updateReactionQuery = `UPDATE reactions SET is_like = ? WHERE  id = ?`

func (p *reactionRepo) UpdateReaction(ctx context.Context, reactionID crypto.UUID, isLike bool) error {
	result, err := p.db.ExecContext(ctx, updateReactionQuery, isLike, reactionID.Value) 	
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

