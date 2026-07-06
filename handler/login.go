package handler

import (
	"encoding/json"
	"net/http"

	"realTime/model"
	"realTime/service"
)

func Login(w http.ResponseWriter, r *http.Request) {
	dataLogin := model.LoginModel{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dataLogin); err != nil {
	http.Error(w, "invalid json", http.StatusBadRequest)
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
