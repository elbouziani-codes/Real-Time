package handler

import (
	"net/http"
)




func NewRouter(authHandler *AuthHandler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/api/register", authHandler.Register)
	router.HandleFunc("/api/login", authHandler.Login)


	return router
}
