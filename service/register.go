package service

import (
	"errors"
	"regexp"
	"strings"

	"realTime/model"
)

var regexpName = regexp.MustCompile("^[A-Za-z]$")
func ValidDataRegister(RegisterForm *model.RegisterModel) error {
	RegisterForm.Email = strings.TrimSpace(RegisterForm.Email)
	RegisterForm.Password = strings.TrimSpace(RegisterForm.Password)
	RegisterForm.NickName = strings.TrimSpace(RegisterForm.NickName)
	RegisterForm.LastName = strings.TrimSpace(RegisterForm.LastName)
	RegisterForm.FirstName = strings.TrimSpace(RegisterForm.FirstName)

	if len(RegisterForm.Password) < 6 || len(RegisterForm.Password) > 20 {
		return errors.New("password length must be between 6 and 20")
	}

	if !emailRegex.MatchString(RegisterForm.Email) {
		if len(RegisterForm.Email) > 40 || len(RegisterForm.Email) < 6 {
			return errors.New("email length must be between 6 and 75")
		}
		return errors.New("is not email")
	} else {
		if !usernameRegex.MatchString(RegisterForm.NickName) {
			return errors.New("invalid nickname")
		}
		if len(RegisterForm.NickName) > 20 || len(RegisterForm.NickName) < 6 {
			return errors.New("nickname length must be between 6 and 20")
		}
	}
	return nil
	return nil
}
