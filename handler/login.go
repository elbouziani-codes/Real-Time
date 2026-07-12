package handler

import (
	"encoding/json"
	"net/http"

	"realTime/domain"
)

type AuthHandler struct {
	svc     domain.AuthService
	userSvc domain.UserService
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var loginForm domain.Credentials

	if err := json.NewDecoder(r.Body).Decode(&loginForm); err != nil {
		// http error handler
	}

	

	// get session
}
