package domain

import "context"


type User struct {
	ID []byte
	NickName string
	LastName string	
	FirstName string
	Email string
	Password string
	Gender string
	Age uint8
	Created_at string
	Updated_at string
}





type UserRepo interface {
	GetByID(context.Context, []byte) (User, error)
	GetByEmail(context.Context, string) (User, error)
	GetByNickName(context.Context, string) (User, error)
	CreateUser(context.Context, User) (User, error)
	GetUsers(context.Context) ([]User, error)
}
