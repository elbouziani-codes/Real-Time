package handler

import (
	"encoding/json"
	"net/http"

	"realTime/model"
	"realTime/service"
)
func Register(w http.ResponseWriter, r *http.Request) {

	dataRegister := model.RegisterModel{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dataRegister); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := service.ValidDataRegister(&dataRegister); err != nil {
		http.Error(w, "error in parser", http.StatusBadRequest)
		return
	}

	newPasswordHash , err := service.PasswordHash(dataRegister.Password)
	if  err != nil {
		http.Error(w, "error in parser", http.StatusBadRequest)
		return
	}
	dataRegister.Password = newPasswordHash
}

