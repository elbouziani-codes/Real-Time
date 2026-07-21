package service

import (
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
)

type UserService struct {
	repo UserRepo
}

type UserRepo interface {
	CreateUser(context.Context, domain.User) error
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) CreateUser(ctx context.Context, user *domain.User) error {

	id, err := crypto.GenerateUUID()
	if err != nil {
		return err
	}

	user.Password, err = crypto.GenerateHash(user.Password)
	if err != nil {
		return err
	}

	user.ID = id
	return svc.repo.CreateUser(ctx, *user)
}
