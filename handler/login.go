package handler

import (
	"net/http"

	"realTime/model"
	"realTime/service"
)

func Login(w http.ResponseWriter, r *http.Request) {
	dataLogin := model.LoginModel{}
	if err := service.ParseJsonLogin(r , &dataLogin); err != nil {
		http.Error(w, "error in parser", http.StatusBadRequest)
		return
	}
	if err := service.ValidDataLogin(&dataLogin); err != nil {
		http.Error(w, "error in parser", http.StatusBadRequest)
		return
	}
	
	if err := service.PasswordHash(&dataLogin); err != nil {
		http.Error(w, "error in parser", http.StatusBadRequest)
		return
	}
}