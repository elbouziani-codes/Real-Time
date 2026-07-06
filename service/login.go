package service

import (
	"errors"
	"regexp"
	"strings"

	"realTime/model"
)

var (
	usernameRegex = regexp.MustCompile(`^[\p{L}\p{N}_]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidDataLogin(loginForm *model.LoginModel) error {
	loginForm.EmailOrNickName = strings.TrimSpace(loginForm.EmailOrNickName)
	loginForm.Password = strings.TrimSpace(loginForm.Password)

	if len(loginForm.Password) < 6 || len(loginForm.Password) > 20 {
		return errors.New("password length must be between 6 and 20")
	}

	if emailRegex.MatchString(loginForm.EmailOrNickName) {
		if len(loginForm.EmailOrNickName) > 40 || len(loginForm.EmailOrNickName) < 6 {
			return errors.New("email length must be between 6 and 75")
		}
	} else {
		if !usernameRegex.MatchString(loginForm.EmailOrNickName) {
			return errors.New("invalid nickname")
		}
		if len(loginForm.EmailOrNickName) > 20 || len(loginForm.EmailOrNickName) < 6 {
			return errors.New("nickname length must be between 6 and 20")
		}
	}
	return nil
}

func PasswordHash(loginForm *model.LoginModel) error {
}
