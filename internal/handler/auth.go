package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"realTime/internal/domain"
	"realTime/internal/service"
)

type AuthRepo interface {
	SaveSession(context.Context, string, string) error
}

type AuthService interface {
	CreateSession(context.Context, string) (string, error)
}

type UserService interface {
	CreateUser(context.Context, service.RegisterInput) (domain.User, error)
}

type AuthHandler struct {
	authSvc AuthService
	userSvc UserService
}

func NewAuthHandler(authSvc AuthService, userSvc UserService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, userSvc: userSvc}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var loginForm domain.Credentials

	if err := json.NewDecoder(r.Body).Decode(&loginForm); err != nil {
		// http error handler
	}

	// get session
}
type output struct{
	NickName string
	Email string
	Age uint8
	Session string

}
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	dataRegister := service.RegisterInput{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dataRegister); err != nil {
		fmt.Println(err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user, err := a.userSvc.CreateUser(r.Context(), dataRegister)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error in parser", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	//fhem
	session, err := a.authSvc.CreateSession(r.Context(), user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "error in parser", http.StatusBadRequest)
		return 
	}
	outputUser := output{NickName : user.NickName,Email : user.Email ,Age:user.Age, Session: session}
	if err := json.NewEncoder(w).Encode(outputUser); err != nil {
		// Avoid writing headers again if Encode fails mid-stream
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
