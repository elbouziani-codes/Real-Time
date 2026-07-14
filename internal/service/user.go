package service

import (
	"context"

	"realTime/crypto"
	"realTime/internal/domain"
)

type UserService struct {
	repo UserRepo
}

type RegisterInput struct {
	NickName  string `json:"nickName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	Age       uint8  `json:"age"`
}

type UserRepo interface {
	GetByID(context.Context, string) (domain.User, error)
	GetByEmail(context.Context, string) (domain.User, error)
	GetByNickName(context.Context, string) (domain.User, error)
	CreateUser(context.Context, domain.User) error
	GetUsers(context.Context, int, int) ([]domain.User, error)
}

func NewUserSevice(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) CreateUser(ctx context.Context, inputs RegisterInput) (domain.User, error) {
	user, err := domain.ValidDataUser(inputs.Email, inputs.Password, inputs.NickName, inputs.LastName, inputs.FirstName, inputs.Gender, inputs.Age)
	if err != nil {
		return domain.User{}, err
	}

	user.ID, err = crypto.GenerateUUID()
	if err != nil {
		return domain.User{}, err
	}

	user.Password, err = crypto.GenerateHash(user.Password)
	if err != nil {
		return domain.User{}, err
	}

	return user, svc.repo.CreateUser(ctx, user)
}
