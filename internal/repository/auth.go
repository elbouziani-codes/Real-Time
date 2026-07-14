package repository

import ( 
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
		&session.ExipireAt)
	if err != nil {
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
		return err
	}
	return nil
}

const DeleteSessionQuery = `
	DELETE FROM sessions WHERE id = ?  
`

func (a *AuthRepo) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := a.db.ExecContext(ctx, DeleteSessionQuery, sessionID)
	if err != nil {
		return err
	}
	return nil
}

const GetSessionByUserIDQuery = `
	SELECT id, user_id FROM sessions WHERE user_id = ?  
`

func (a *AuthRepo) GetBySessionID(ctx context.Context, UserID string) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, GetSessionSessionByIDQuery, userID)
	return scanSession(row)
}

const GetSessionByIDQuery = `
	SELECT id, user_id FROM sessions WHERE id = ?  
`

func (u *AuthRepo) GetByID(ctx context.Context, sessionID string) (domain.Session, error) {
	row := u.db.QueryRowContext(ctx, GetSessionByIDQuery, sessionID)
	return scanSession(row)
}
// getByUser

// getById
