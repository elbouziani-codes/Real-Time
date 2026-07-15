package handler

import (
	"errors"
	"context"
	"fmt"
	"encoding/json"
	"net/http"
	"realTime/internal/domain"
	"realTime/internal/service"
)

type AuthRepo interface {
	SaveSession(context.Context, string, string) error
}

type AuthService interface {
	Login(context.Context, domain.Credentials, int) (string, error)
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

type output struct{
	NickName string
	Email string
	Age int
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
	// ERROR HANDLING SECTION
	if err != nil {
		var valErr domain.Error      
		
		if errors.As(err, &valErr) {
			fmt.Println(err)
			switch valErr.Code {
				case domain.ConflictCode:  
					http.Error(w, err.Error(), http.StatusConflict)
					return
				case domain.BadFormatCode: 
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				default: 
					http.Error(w, "InternalServerError", http.StatusInternalServerError)
					return
			}			
		 }
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return 

	}
	//////

	

	session, err := a.authSvc.CreateSession(r.Context(), user.ID)
	if err != nil {	
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return 
	}

	w.Header().Set("Content-Type", "application/json")

	outputUser := output{NickName : user.NickName,Email : user.Email ,Age:user.Age, Session: session}
	if err := json.NewEncoder(w).Encode(outputUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds domain.Credentials 
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&creds); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	idType, err := domain.LoginValidation(creds)  
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//instead of repeating code I will add helper for that
		return
	}
 	sessionID, err := a.authSvc.Login(r.Context(), creds, idType)
	if err := json.NewEncoder(w).Encode(sessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func (a *AuthHandler) Error(err error, w http.ResponseWriter, r *http.Request) {
		if errors.As(err, &valErr) {
			fmt.Println(err)
			switch valErr.Code {
				case domain.ConflictCode:  
					http.Error(w, err.Error(), http.StatusConflict)
					return
				case domain.BadFormatCode: 
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				default: 
					http.Error(w, "InternalServerError", http.StatusInternalServerError)
					return
			}			
		}	
}


