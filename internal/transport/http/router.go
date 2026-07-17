package handler

import (
	"net/http"
)



type middleWare interface {
	Auth(http.Handler) http.Handler
}

func NewRouter(authHandler *AuthHandler, testHandler TestHandler, middleware middleWare) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/api/register", authHandler.Register)
	router.HandleFunc("/api/login", authHandler.Login)
	router.Handle("/api/test", middleware.Auth(http.HandlerFunc(testHandler.Test)))
	return router
}
