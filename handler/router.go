package handler

import (
	"net/http"

	"realTime/model"
)

type servece interface {
	ValidDataLogin(loginForm *model.LoginModel) error
	ValidDataRegister(RegisterForm *model.RegisterModel) error
}

type Handler struct{
	sr servece
}

func Loadhandler(serveces servece) *Handler {
	return &Handler{
		sr: serveces,
	}
}
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/Register", Register)
	mux.HandleFunc("/api/Login", Login)
	mux.HandleFunc("/api/Logout", Logout)
}
