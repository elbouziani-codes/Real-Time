package service

import (
	"context"
	"errors"
	"realTime/crypto"
	"realTime/internal/domain"
	"time"
)

type AuthService struct {
	authRepo   authRepository
	userGetter userGetterRepo
}

type authRepository interface {
	DeleteSession(context.Context, crypto.UUID) error
	SaveSession(context.Context, crypto.UUID, crypto.UUID) error
	GetByID(context.Context, crypto.UUID) (domain.Session, error)
}

type userGetterRepo interface {
	GetByNickName(context.Context, string) (domain.User, error)
	GetByEmail(context.Context, string) (domain.User, error)
}

func NewAuthService(authRepo authRepository, userGetter userGetterRepo) *AuthService {
	return &AuthService{authRepo: authRepo, userGetter: userGetter}
}

func (a *AuthService) CreateSession(ctx context.Context, userID crypto.UUID) (crypto.UUID, error) {
	id, err := crypto.GenerateUUID()
	if err != nil {
		return crypto.Nil, err
	}

	err = a.authRepo.SaveSession(ctx, id, userID)
	if err != nil {
		return crypto.Nil, err
	}

	return id, nil
}

func (a *AuthService) ValueidateSession(ctx context.Context, sessionID crypto.UUID) (domain.Session, error) {
	session, err := a.authRepo.GetByID(ctx, sessionID)
	if err != nil {
		return domain.Session{}, err
	}
	if time.Now().Unix() > session.ExpireAt {
		err := a.Logout(ctx, session.UserID)
		if err != nil {
			return domain.Session{}, err
		}
		return domain.Session{}, domain.Error{Message: "invalid credentials", Code: domain.UnauthorizedCode}
	}
	return session, nil
}

func (a *AuthService) Logout(ctx context.Context, userID crypto.UUID) error {
	err := a.authRepo.DeleteSession(ctx, userID)
	if err != nil {
		var buckErr domain.Error
		if errors.As(err, buckErr) {
			if buckErr.Code == domain.NotFoundCode {
				return nil
			}
		}
		return err
	}
	return nil
}

func (a *AuthService) Login(ctx context.Context, creds domain.Credentials, idType int) (crypto.UUID, error) {
	// phas_1 finding target user
	var userID crypto.UUID
	var hashedPassword string
	if idType == domain.EmailType {
		user, err := a.userGetter.GetByEmail(ctx, creds.Identifier)
		if err != nil {
			return crypto.Nil, domain.Error{Message: "invalid credentials", Code: domain.UnauthorizedCode}

		}
		userID = user.ID
		hashedPassword = user.Password
	} else {
		user, err := a.userGetter.GetByNickName(ctx, creds.Identifier)
		if err != nil {
			return crypto.Nil, domain.Error{Message: "invalid credentials", Code: domain.UnauthorizedCode}

		}
		userID = user.ID
		hashedPassword = user.Password
	}
	// phase_2 matching password
	err := crypto.CompareHashWithPassword(hashedPassword, creds.Password)
	if err != nil {
		return crypto.Nil, domain.Error{Message: "invalid credentials", Code: domain.UnauthorizedCode}
	}

	// clearing old session
	err = a.Logout(ctx, userID)
	if err != nil {
		return crypto.Nil, err
	}

	//spawnning new one
	sessionID, err := a.CreateSession(ctx, userID)
	if err != nil {
		return crypto.Nil, err
	}
	return sessionID, nil
}
