package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"realTime/internal/repository"
	"realTime/internal/service"
	handler "realTime/internal/transport/http"
	ws "realTime/internal/transport/webSocket"

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
	chatRepo := repository.NewChatRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	reactionRepo := repository.NewReactionRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)




	// new service
	authSvc := service.NewAuthService(authRepo, userRepo)
	userSvc := service.NewUserService(userRepo)
	postSvc := service.NewPostService(postRepo)
	charSvc := service.NewChatService(chatRepo, userRepo)
	commentSvc := service.NewCommentService(commentRepo)
	reactionSvc := service.NewReactionService(reactionRepo)
	categorySvc := service.NewCategoryService(categoryRepo)



	// new handler
	authHandler := handler.NewAuthHandler(authSvc, userSvc)
	postHandler := handler.NewPostHandler(postSvc)
	wsHandler := ws.NewHandleWs(charSvc, userSvc)
	chatHandler := handler.NewHandleChat(charSvc, userSvc)
	commentHandler := handler.NewCommentHandler(commentSvc, postSvc)
	reactionHandler := handler.NewReactionHandler(reactionSvc, postSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)

	middleware := middleware.NewMiddleware(authSvc)
	router := handler.NewRouter(authHandler, postHandler, commentHandler, reactionHandler, wsHandler, chatHandler, categoryHandler, middleware)

	spa := spaHandler()
	log.Println("Database Initialised")
	if err := http.ListenAndServe(conf.PortMux2, spa(router)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}



func spaHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fileServer := http.FileServer(http.Dir("./web"))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			if strings.HasPrefix(r.URL.Path ,"/api/"){
				next.ServeHTTP(w, r)
				return 
			}

			path := "./web" + r.URL.Path
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				http.ServeFile(
					w,
					r,
					"./web/index.html",
				)
				return
			}
			// Serve frontend files
			fileServer.ServeHTTP(w,r)
		})
	}
}