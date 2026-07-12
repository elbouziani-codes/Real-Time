package service

import (
	"context"

	"realTime/crypto"
	"realTime/internal/domain"
)

type AuthService struct {
	authRepo AuthRepo
	userRepo UserRepo
}
type AuthRepo interface {
	SaveSession(context.Context, string, string) error
}
type RegisterRepo interface {
	GetByID(context.Context, string) (domain.User, error)
}

func NewAuthService(authRepo AuthRepo, UserRepo UserRepo) *AuthService {
	return &AuthService{authRepo: authRepo, userRepo: UserRepo}
}

func (a AuthService) CreateSession(ctx context.Context, userID string) (string, error) {
	id, err := crypto.GenerateUUID()
	if err != nil {
		return "", err
	}

	err = a.authRepo.SaveSession(ctx, userID, id)
	if err != nil {
		return "", err
	}

	return id, nil
}
