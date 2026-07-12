package domain

import (
	"context"
	"errors"
	"regexp"
	"strings"
)

type User struct {
	ID         string
	NickName   string
	LastName   string
	FirstName  string
	Email      string
	Password   string
	Gender     string
	Age        uint8
	Created_at string
	Updated_at string
}



type UserRepo interface {
	GetByID(context.Context, string) (User, error)
	GetByEmail(context.Context, string) (User, error)
	GetByNickName(context.Context, string) (User, error)
	CreateUser(context.Context, User) (User, error)
	GetUsers(context.Context) ([]User, error)
}


type UserService interface {
	CreateUser(context.Context, string) (string, error)
}

var nameRegex = regexp.MustCompile(`^[\p{L}]+(?:[-'][\p{L}]+)*(?: [\p{L}]+(?:[-'][\p{L}]+)*)*$`)

func ValidDataUser(email, password, nickName, lastName, firstName, gender string, age uint8) (User, error) {
	var user User
	user.Email = strings.TrimSpace(email)
	user.Password = strings.TrimSpace(password)
	user.NickName = strings.TrimSpace(nickName)
	user.LastName = strings.TrimSpace(lastName)
	user.FirstName = strings.TrimSpace(firstName)
	user.Gender = strings.TrimSpace(gender)
	user.Age = age
	if len(user.Password) < 6 || len(user.Password) > 20 {
		return User{}, errors.New("password length must be between 6 and 20")
	}

	if len(user.Email) < 6 || len(user.Email) > 75 {
		return User{}, errors.New("email length must be between 6 and 75")
	}

	if !emailRegex.MatchString(user.Email) {
		return User{}, errors.New("invalid email")
	}

	if !usernameRegex.MatchString(user.NickName) {
		return User{}, errors.New("invalid nickname")
	}

	if len(user.NickName) > 20 || len(user.NickName) < 6 {
		return User{}, errors.New("nickname length must be between 6 and 20")
	}

	if !nameRegex.MatchString(user.FirstName) || !nameRegex.MatchString(user.LastName) {
		return User{}, errors.New("invalid FirstName or LastName")
	}

	if (len(user.FirstName) > 25 || len(user.FirstName) < 4) || (len(user.LastName) > 25 || len(user.LastName) < 4) {
		return User{}, errors.New("FirstName and LastName length must be between 4 and 25 ")
	}
	if user.Gender != "Man" && user.Gender != "Woman" {
		return User{}, errors.New("Value Gender is not valid")
	}
	return user, nil
}
