package main

import (
	"log"
	"net/http"
	"os"

	"realTime/internal/repository"
	"realTime/internal/service"
	"realTime/internal/transport/http"

	"realTime/internal/config"
	"realTime/internal/repository/sqlite"
	"realTime/internal/transport/http/middleware"
)

func main() {
	conf := config.Load()

	db, err := sqlite.Open(conf.DBPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// new repo
	authRepo := repository.NewAuthRepo(db)
	userRepo := repository.NewUserRepo(db)
	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	reactionRepo := repository.NewReactionRepo(db)



	// new service
	authSvc := service.NewAuthService(authRepo, userRepo)
	userSvc := service.NewUserService(userRepo)
	postSvc := service.NewPostService(postRepo)
	commentSvc := service.NewCommentService(commentRepo)
	reactionSvc := service.NewReactionService(reactionRepo)



	// new handler
	authHandler := handler.NewAuthHandler(authSvc, userSvc)
	postHandler := handler.NewPostHandler(postSvc)
	commentHandler := handler.NewCommentHandler(commentSvc, postSvc)
	reactionHandler := handler.NewReactionHandler(reactionSvc, postSvc)


	middleware := middleware.NewMiddleware(authSvc)
	router := handler.NewRouter(authHandler, postHandler, commentHandler, reactionHandler, middleware)
	http.ListenAndServe(":8081", router)
	log.Println("Database Initialised")
}
