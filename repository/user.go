package repository

import (
	"context"
	"realTime/domain"
)

type UserRepo struct {
	db DBTX
}

func NewUserRepo(db DBTX) *UserRepo {
	return &UserRepo{db: db}
}
/* 
type User struct {
	ID []byte
	NickName string
	LastName string	
	FirstName string
	Email string
	Gender string
	Age uint8
	Created_at string
	Updated_at string
}


*/


const CreateUserQuery = `INSERT INTO users 
(email, password_hash, nick_name, last_name, first_name, age, gender) VALUES 
(?, ?, ?, ?, ?, ?, ?) `

func (u *UserRepo) CreateUser(ctx context.Context, user domain.User) error {
	row, err := u.db.ExecContext(ctx, CreateUserQuery, 
		user.Email,
		user.Password,
		user.NickName,
		user.LastName,
		user.FirstName,
		user.Age,
		user.Gender) 
	if err != nil{
		return err
	}
	_, err  = row.LastInsertId()

	if err != nil{
		return err
	}
	return nil 
}