package handler

import (
	"net/http"
)


func RegisterHandlers(mux *http.ServeMux){
	mux.HandleFunc("/api/Register",Register)
	mux.HandleFunc("/api/Login",Login)
}