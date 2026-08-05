package domain

import (
	"realTime/crypto"
	"regexp"
	"strings"
)

type User struct {
	ID        crypto.UUID
	NickName  string
	LastName  string
	FirstName string
	Email     string
	Password  string
	Gender    string
	Age       int
	CreatedAt int
	UpdatedAt int
}

type UserProfile struct {
	ID        crypto.UUID
	NickName  string
	LastName  string
	FirstName string
	Gender    string
	Age       int
	CreatedAt int
}
type RegisterRequest struct {
	NickName  string `json:"nick_name"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
}

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]*$`) // would deleted later

func ValueidateUserInfo(registerRequest RegisterRequest) (User, error) {
	var user User
	user.Email = strings.TrimSpace(registerRequest.Email)
	user.Password = registerRequest.Password
	user.NickName = strings.TrimSpace(registerRequest.NickName)
	user.LastName = strings.TrimSpace(registerRequest.LastName)
	user.FirstName = strings.TrimSpace(registerRequest.FirstName)
	user.Gender = strings.TrimSpace(registerRequest.Gender)
	user.Age = registerRequest.Age
	if len(user.Password) < 8 || len(user.Password) > 20 {
		return user, Error{Message: "password length must be between 8 and 20", Code: BadFormatCode}
	}

	if len(user.Email) < 6 || len(user.Email) > 75 {
		return user, Error{Message: "email length must be between 6 and 75", Code: BadFormatCode} // 400
	}

	if !emailRegex.MatchString(user.Email) {
		return user, Error{Message: "email invalid format", Code: BadFormatCode}
	}

	if !usernameRegex.MatchString(user.NickName) {
		return user, Error{Message: "nickname invalid format", Code: BadFormatCode}
	}

	if len(user.NickName) > 20 || len(user.NickName) < 2 {
		return user, Error{Message: "nickname length must be between 2 and 20", Code: BadFormatCode}
	}

	if !nameRegex.MatchString(user.FirstName) || !nameRegex.MatchString(user.LastName) {
		return user, Error{Message: "invalid FirstName or LastName", Code: BadFormatCode}
	}

	if (len(user.FirstName) > 25 || len(user.FirstName) < 2) || (len(user.LastName) > 25 || len(user.LastName) < 2) {
		return user, Error{Message: "firstName and LastName length must be between 2 and 25", Code: BadFormatCode}
	}
	if user.Gender != "man" && user.Gender != "woman" {
		return user, Error{Message: "value gender is not valid", Code: BadFormatCode}
	}
	if user.Age < 14 {
		return user, Error{Message: "you are too old", Code: BadFormatCode}
	}
	if user.Age > 200 {
		return user, Error{Message: "invalid age", Code: BadFormatCode}
	}
	return user, nil
}
