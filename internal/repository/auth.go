package repository

import (
	"context"
	"fmt"
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
		fmt.Println(err)
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
	fmt.Println(session)
	return session, nil
}

const createSessionQuery = `
	INSERT into sessions (id, user_id) VALUES (?, ?)
`

func (a *AuthRepo) SaveSession(ctx context.Context, sessionID, userID crypto.UUID) error {
	fmt.Println(userID, sessionID)
	_, err := a.db.ExecContext(ctx, createSessionQuery, sessionID.Value, userID.Value)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const DeleteSessionQuery = `
	DELETE FROM sessions WHERE user_id = ?  
`

func (a *AuthRepo) DeleteSession(ctx context.Context, userID crypto.UUID) error {
	_, err := a.db.ExecContext(ctx, DeleteSessionQuery, userID.Value)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const GetSessionByUserIDQuery = `
	SELECT id, user_id FROM sessions WHERE user_id = ?  
`

func (a *AuthRepo) GetByUserID(ctx context.Context, userID crypto.UUID) (domain.Session, error) {
	row := a.db.QueryRowContext(ctx, GetSessionByUserIDQuery, userID.Value)
	return scanSession(row)
}

const GetSessionByIDQuery = `
	SELECT id, user_id, created_at, expire_at FROM sessions WHERE id = ?  
`

func (a *AuthRepo) GetByID(ctx context.Context, sessionID crypto.UUID) (domain.Session, error) {
	row := a.db.QueryRowContext(ctx, GetSessionByIDQuery, sessionID.Value)
	return scanSession(row)
}
