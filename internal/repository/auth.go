package repository

import "context"

type AuthRepo struct {
	db DBTX
}

func NewAuthRepo(db DBTX) *AuthRepo {
	return &AuthRepo{db: db}
}

const createSessionQuery = `
	INSERT into sessions (user_id, id) VALUES (?, ?)
`

func (a *AuthRepo) SaveSession(ctx context.Context, UserID, sessionID string) error {
	_, err := a.db.ExecContext(ctx, createSessionQuery, UserID, sessionID)
	if err != nil {
		return err
	}
	return nil
}

const DeleteSessionQuery = `
	DELETE into sessions (user_id, id) VALUES (?, ?)
`

// delete

// getByUser

// getById
