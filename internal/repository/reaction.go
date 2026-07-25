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


const getReactionsQuery = `
SELECT C.id,  C.content, C.created_at, C.parent_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at
FROM reactions C  
JOIN users U 
ON C.author_id = U.id
WHERE C.parent_id = ?
`

func scanReaction(row scanner) (*domain.ReactionInfo, error) {
	reaction := domain.ReactionInfo{}
	err := row.Scan(
		&reaction.ID.Value,
		&reaction.IsLike,
		&reaction.CreatedAt,
		&reaction.ParentID.Value,
		&reaction.Author.ID.Value,
		&reaction.Author.NickName,
		&reaction.Author.LastName,
		&reaction.Author.FirstName,
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
	_, err := p.db.ExecContext(ctx, saveReactionQuery, reaction.ID.Value, reaction.ParentID.Value, reaction.AuthorID.Value, reaction.Content)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const getReactionQuery = `
SELECT C.id,  C.content, C.created_at, C.parent_id, U.id, U.nick_name, U.first_name, U.last_name, U.gender, U.age, U.created_at
FROM reactions C  
JOIN users U 
ON C.author_id = U.id
WHERE C.id = ?
`


func (p *reactionRepo) GetReaction(ctx context.Context, reactionID crypto.UUID) (*domain.ReactionInfo, error) {
	row := p.db.QueryRowContext(ctx, getReactionQuery, reactionID.Value)	
	return scanReaction(row) 
}


//const getReactionsQuery = `SELECT id, author_id, title, content, created_at, updated_at FROM reactions 
//						ORDER BY created_at DESC 
//						LIMIT ? OFFSET ?; `

/*func (p *reactionRepo) GetReactions(ctx context.Context, parentID crypto.UUID) ([]*domain.ReactionInfo, error) {
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
		fmt.Println(reactions)
		reactions = append(reactions, reaction)
	}
	return reactions, nil 
}*/

const deleteReactionQuery = `DELETE FROM reactions WHERE id = ?`

func (p *reactionRepo) DeleteReaction(ctx context.Context, reactionID crypto.UUID) error {
	_, err := p.db.ExecContext(ctx, deleteReactionQuery, reactionID.Value)
	if err != nil {
		return err
	}
	return nil
}

func (p *reactionRepo) ExecPatchQuery(ctx context.Context, query string, reactionID crypto.UUID, args []any) error {
	result, err := p.db.ExecContext(ctx, query, append(args, reactionID.Value)...) 	
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

