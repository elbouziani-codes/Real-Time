package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"realTime/model"
)
var usernameRegex = regexp.MustCompile(`^[\p{L}\p{N}_]+$`)
func ParseJsonLogin(r *http.Request, loginForm *model.LoginModel) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(loginForm)
}

func ValidDataLogin(loginForm *model.LoginModel) error {
	loginForm.EmailOrNickName = strings.TrimSpace(loginForm.EmailOrNickName)
	loginForm.Password = strings.TrimSpace(loginForm.Password)

	if len(loginForm.Password) < 6 || len(loginForm.Password) > 20 {
		return errors.New("password length must be between 6 and 20")
	}

	if _, err := mail.ParseAddress(loginForm.EmailOrNickName); err == nil {
		if len(loginForm.EmailOrNickName) > 75 || len(loginForm.EmailOrNickName) < 6 {
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

func PasswordHash(loginForm *model.LoginModel)error{

}