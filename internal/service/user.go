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
	GetByID(context.Context, crypto.UUID) (domain.User, error)
	CreateUser(context.Context, domain.User) error
	GetUserProfile(context.Context, crypto.UUID, crypto.UUID) (*domain.UserProfile, error)
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) GetUser(ctx context.Context, userID, requesterID crypto.UUID) (*domain.UserProfile, error) {
	return svc.repo.GetUserProfile(ctx, userID, requesterID) 	
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


func (svc *UserService) GetByID(ctx context.Context, id crypto.UUID) (domain.User , error){
	return svc.repo.GetByID(ctx, id)
}
