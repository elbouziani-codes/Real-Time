package domain

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"
)



type Credentials struct {
	EmailOrNickName string `json:"identifier"`
	Password        string `json:"password"`
}


var (
	usernameRegex = regexp.MustCompile(`^[\p{L}\p{N}_]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)


func NewLoginForm(requestBody io.ReadCloser) (LoginModel, error) {
	var loginForm LoginModel
	if err := json.NewDecoder(requestBody).Decode(&loginForm); err != nil {
		return LoginModel{}, err
	}

	loginForm.EmailOrNickName = strings.TrimSpace(loginForm.EmailOrNickName)
	loginForm.Password = strings.TrimSpace(loginForm.Password)

	if len(loginForm.Password) < 6 || len(loginForm.Password) > 20 {
		return LoginModel{},errors.New("password length must be between 6 and 20")
	}

	if emailRegex.MatchString(loginForm.EmailOrNickName) {
		if len(loginForm.EmailOrNickName) > 50 || len(loginForm.EmailOrNickName) < 6 {
			return LoginModel{},errors.New("email length must be between 6 and 50")
		}
	} else {
		if !usernameRegex.MatchString(loginForm.EmailOrNickName) {
			return LoginModel{},errors.New("invalid nickname")
		}
		if len(loginForm.EmailOrNickName) > 20 || len(loginForm.EmailOrNickName) < 6 {
			return LoginModel{},errors.New("nickname length must be between 6 and 20")
		}
	}

	return loginForm , nil
}
