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
	Age        int
	CreatedAt  int 
	UpdatedAt  int 
}




var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9]{2,40}$`) // would deleted later

func ValidDataUser(email, password, nickName, lastName, firstName, gender string, age int) (User, error) {
	var user User
	user.Email = strings.TrimSpace(email)
	user.Password = password
	user.NickName = strings.TrimSpace(nickName)
	user.LastName = strings.TrimSpace(lastName)
	user.FirstName = strings.TrimSpace(firstName)
	user.Gender = strings.TrimSpace(gender)
	user.Age = age
	if len(user.Password) < 6 || len(user.Password) > 20 {
			return User{}, Error{Field: "password", Message: "length must be between 6 and 20", Code: BadFormatCode}
	}

	if len(user.Email) < 6 || len(user.Email) > 75 {
			return User{}, Error{Field: "email", Message: "length must be between 6 and 75", Code: BadFormatCode} // 400 
	}

	if !emailRegex.MatchString(user.Email) {
		return User{}, Error{Field: "email", Message: "invalid format", Code: BadFormatCode} 
	}

	if !usernameRegex.MatchString(user.NickName) {
		return User{}, Error{Field: "nickname", Message: "invalid format", Code: BadFormatCode}  
	}

	if len(user.NickName) > 20 || len(user.NickName) < 6 {
		return User{}, Error{Field: "nickname", Message: "nickname length must be between 6 and 20", Code: BadFormatCode}
	}

	if !nameRegex.MatchString(user.FirstName) || !nameRegex.MatchString(user.LastName) {
		return User{}, Error{Field: "fullname", Message: "invalid FirstName or LastName", Code: BadFormatCode}
	}

	if (len(user.FirstName) > 25 || len(user.FirstName) < 4) || (len(user.LastName) > 25 || len(user.LastName) < 4) {
		return User{}, Error{Field: "fullname", Message: "FirstName and LastName length must be between 4 and 25", Code: BadFormatCode}
	}
	if user.Gender != "Man" && user.Gender != "Woman" {
			return User{}, Error{Field: "gender", Message: "Value Gender is not valid", Code: BadFormatCode}
	}
	return user, nil
}
