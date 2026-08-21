package repository

import (
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
)

type AuthRepo struct {
	db DBTX
}

func NewAuthRepo(db DBTX) *AuthRepo {
	return &AuthRepo{db: db}
}

func scanSession(row scanner) (domain.Session, error) {
	session := domain.Session{}
	//	session.ID = &crypto.Nil
	//	session.UserID = &crypto.Nil

	//	var sessionID any
	//	var userID any
	err := row.Scan(
		&session.ID.Value,
		&session.UserID.Value,
		&session.CreatedAt,
		&session.ExpireAt)
	if err != nil {
		return domain.Session{}, sqlite.TranslateError(err)
	}
	/*
		err = session.ID.Scan(sessionID)
		if err != nil {
			return session, err

		}

		err = session.UserID.Scan(userID)
		if err != nil {
			return session, err

		}*/
	return session, nil
}

const createSessionQuery = `
	INSERT into sessions (id, user_id) VALUES (?, ?)
`

func (a *AuthRepo) SaveSession(ctx context.Context, sessionID, userID crypto.UUID) error {
	_, err := a.db.ExecContext(ctx, createSessionQuery, sessionID.Value, userID.Value)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const DeleteSessionQuery = `
	DELETE FROM sessions WHERE user_id = ?  
`

// DeleteSession reports NotFound when the user had no session, so callers can
// tell "nothing to clear" apart from a genuine failure.
func (a *AuthRepo) DeleteSession(ctx context.Context, userID crypto.UUID) error {
	result, err := a.db.ExecContext(ctx, DeleteSessionQuery, userID.Value)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return sqlite.TranslateError(err)
	}
	if rows == 0 {
		return domain.Error{Message: "not found", Code: domain.NotFoundCode}
	}
	return nil
}

const GetSessionByIDQuery = `
	SELECT id, user_id, created_at, expire_at FROM sessions WHERE id = ?  
`

func (a *AuthRepo) GetByID(ctx context.Context, sessionID crypto.UUID) (domain.Session, error) {
	row := a.db.QueryRowContext(ctx, GetSessionByIDQuery, sessionID.Value)
	return scanSession(row)
}
