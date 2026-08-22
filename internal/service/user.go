package service

import (
	"context"

	"realTime/crypto"
	"realTime/internal/domain"
	"uuid"
)

type UserService struct {
	repo UserRepo
}

type UserRepo interface {
	GetByID(context.Context, uuid.UUID) (domain.User, error)
	CreateUser(context.Context, domain.User) error
	GetUserProfile(context.Context, uuid.UUID, uuid.UUID) (*domain.UserProfile, error)
	GetUsers(context.Context, uuid.UUID, int, uuid.UUID) ([]domain.UserContact, error)
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) GetUser(ctx context.Context, userID, requesterID uuid.UUID) (*domain.UserProfile, error) {
	return svc.repo.GetUserProfile(ctx, userID, requesterID)
}

// GetUsers lists the people visible to userID, most recently talked with first.
func (svc *UserService) GetUsers(ctx context.Context, userID uuid.UUID, limit int, cursor uuid.UUID) ([]domain.UserContact, error) {
	return svc.repo.GetUsers(ctx, userID, limit, cursor)
}

func (svc *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	user.ID = uuid.NewV4()

	hashed, err := crypto.GenerateHash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return svc.repo.CreateUser(ctx, *user)
}

func (svc *UserService) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return svc.repo.GetByID(ctx, id)
}
