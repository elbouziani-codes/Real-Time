package domain

import (
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




var nameRegex = regexp.MustCompile(`^[\p{L}]+(?:[-'][\p{L}]+)*(?: [\p{L}]+(?:[-'][\p{L}]+)*)*$`) // would deleted later

func ValidDataUser(email, password, nickName, lastName, firstName, gender string, age uint8) (User, error) {
	var user User
	user.Email = strings.TrimSpace(email)
	user.Password = password
	user.NickName = strings.TrimSpace(nickName)
	user.LastName = strings.TrimSpace(lastName)
	user.FirstName = strings.TrimSpace(firstName)
	user.Gender = strings.TrimSpace(gender)
	user.Age = age
	if len(user.Password) < 6 || len(user.Password) > 20 {
			return User{}, &ValidationError{Field: "password", Message: "length must be between 6 and 20", Code: 400}
	}

	if len(user.Email) < 6 || len(user.Email) > 75 {
			return User{}, &ValidationError{Field: "email", Message: "length must be between 6 and 75", Code: 400} // 400 
	}

	if !emailRegex.MatchString(user.Email) {
		return User{}, &ValidationError{Field: "email", Message: "invalid format", Code: 400} 
	}

	if !usernameRegex.MatchString(user.NickName) {
		return User{}, &ValidationError{Field: "nickname", Message: "invalid format", Code: 400}  
	}

	if len(user.NickName) > 20 || len(user.NickName) < 6 {
		return User{}, &ValidationError{Field: "nickname", Message: "nickname length must be between 6 and 20", Code: 400}
	}

	if !nameRegex.MatchString(user.FirstName) || !nameRegex.MatchString(user.LastName) {
		return User{}, &ValidationError{Field: "fullname", Message: "invalid FirstName or LastName", Code: 400}
	}

	if (len(user.FirstName) > 25 || len(user.FirstName) < 4) || (len(user.LastName) > 25 || len(user.LastName) < 4) {
		return User{}, &ValidationError{Field: "fullname", Message: "FirstName and LastName length must be between 4 and 25", Code: 400}
	}
	if user.Gender != "Man" && user.Gender != "Woman" {
			return User{}, &ValidationError{Field: "gender", Message: "Value Gender is not valid", Code: 400}
	}
	return user, nil
}
