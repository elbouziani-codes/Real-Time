package service

import (
	"errors"
	"regexp"
	"strings"

	"realTime/model"

	"golang.org/x/crypto/bcrypt"
)

var nameRegex = regexp.MustCompile(`^[\p{L}]+(?:[-'][\p{L}]+)*(?: [\p{L}]+(?:[-'][\p{L}]+)*)*$`)

func ValidDataRegister(RegisterForm *model.RegisterModel) error {
	RegisterForm.Email = strings.TrimSpace(RegisterForm.Email)
	RegisterForm.Password = strings.TrimSpace(RegisterForm.Password)
	RegisterForm.NickName = strings.TrimSpace(RegisterForm.NickName)
	RegisterForm.LastName = strings.TrimSpace(RegisterForm.LastName)
	RegisterForm.FirstName = strings.TrimSpace(RegisterForm.FirstName)
	//Gender
	//Age 

	if len(RegisterForm.Password) < 6 || len(RegisterForm.Password) > 20 {
		return errors.New("password length must be between 6 and 20")
	}

	if len(RegisterForm.Email) < 6 || len(RegisterForm.Email) > 75 {
		return errors.New("email length must be between 6 and 75")
	}

	if !emailRegex.MatchString(RegisterForm.Email) {
		return errors.New("invalid email")
	}

	if !usernameRegex.MatchString(RegisterForm.NickName) {
		return errors.New("invalid nickname")
	}

	if len(RegisterForm.NickName) > 20 || len(RegisterForm.NickName) < 6 {
		return errors.New("nickname length must be between 6 and 20")
	}

	if !nameRegex.MatchString(RegisterForm.FirstName) || !nameRegex.MatchString(RegisterForm.LastName) {
		return errors.New("invalid FirstName or LastName")
	}

	if (len(RegisterForm.FirstName) > 25 || len(RegisterForm.FirstName) < 4) || (len(RegisterForm.LastName) > 25 || len(RegisterForm.LastName) < 4) {
		return errors.New("FirstName and LastName length must be between 4 and 25 ")
	}

	return nil
}

func (S Servece) PasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
