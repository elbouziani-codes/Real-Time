package main

import (
	"log"
	"net/http"
	"os"

	"realTime/internal/database"
	"realTime/internal/transport/http"
	"realTime/internal/repository"
	"realTime/internal/service"

	"realTime/internal/config"
	"realTime/internal/transport/http/middleware"

)

func main() {
	conf := config.Load()

	db, err := database.Open(conf.DBPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// new repo
	authRepo := repository.NewAuthRepo(db)
	userRepo := repository.NewUserRepo(db)

	// new service
	authSvc := service.NewAuthService(authRepo, userRepo)
	userSvc := service.NewUserSevice(userRepo)

	// new handler

	authHandler := handler.NewAuthHandler(authSvc, userSvc)
	testHandler := handler.NewTestHandler(authSvc, userSvc)
		
	middleware := middleware.NewMiddleware(authSvc)
	router := handler.NewRouter(authHandler, testHandler, middleware)
	http.ListenAndServe(":8080", router)
	log.Println("Database Initialised")
}
