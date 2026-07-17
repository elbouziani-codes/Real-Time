package repository

import ( 
		"fmt"
		"context"
		"realTime/internal/domain"
)

type AuthRepo struct {
	db DBTX
}

func NewAuthRepo(db DBTX) *AuthRepo {
	return &AuthRepo{db: db}
}

func scanSession(row scanner) (domain.Session, error) {
	session := domain.Session{}
	err := row.Scan(	
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpireAt)
	if err != nil {
		fmt.Println(err)
		return domain.Session{},  TranslateError(err)
	}
	return session, nil
}


const createSessionQuery = `
	INSERT into sessions (user_id, id) VALUES (?, ?)
`

func (a *AuthRepo) SaveSession(ctx context.Context, UserID, sessionID string) error {
	_, err := a.db.ExecContext(ctx, createSessionQuery, UserID, sessionID)
	if err != nil {
		return  TranslateError(err)
	}
	return nil
}

const DeleteSessionQuery = `
	DELETE FROM sessions WHERE user_id = ?  
`

func (a *AuthRepo) DeleteSession(ctx context.Context, userID string) error {
	_, err := a.db.ExecContext(ctx, DeleteSessionQuery, userID)
	if err != nil {
		return  TranslateError(err)
	}
	return nil
}

const GetSessionByUserIDQuery = `
	SELECT id, user_id FROM sessions WHERE user_id = ?  
`

func (a *AuthRepo) GetByUserID(ctx context.Context, userID string) (domain.Session, error) {
	row := a.db.QueryRowContext(ctx, GetSessionByUserIDQuery, userID)
	return scanSession(row)
}

const GetSessionByIDQuery = `
	SELECT id, user_id, created_at, expire_at FROM sessions WHERE id = ?  
`

func (a *AuthRepo) GetByID(ctx context.Context, sessionID string) (domain.Session, error) {
	row := a.db.QueryRowContext(ctx, GetSessionByIDQuery,  sessionID)
	return scanSession(row)
}

