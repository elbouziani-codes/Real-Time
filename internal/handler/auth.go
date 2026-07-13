package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"
	"errors"
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
	// TODO LATER 
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
	decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&dataRegister); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
		}
	

	user, err := a.userSvc.CreateUser(r.Context(), dataRegister)
	if err != nil {
		var valErr domain.ValidationError
		
		if errors.As(err, &valErr) {
			switch valErr.Code {
				case 409:  
					http.Error(w, err.Error(), http.StatusBadRequest)
				case 400: 
					http.Error(w, err.Error(), http.StatusBadRequest)
				default: 
					http.Error(w, "InternalServerError", http.StatusInternalServerError)
					return
			}			
		 }
		fmt.Println(err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return 

	}
	w.Header().Set("Content-Type", "application/json")

	//fhem
	session, err := a.authSvc.CreateSession(r.Context(), user.ID)
	if err != nil {	
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return 
	}
	outputUser := output{NickName : user.NickName,Email : user.Email ,Age:user.Age, Session: session}
	if err := json.NewEncoder(w).Encode(outputUser); err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
