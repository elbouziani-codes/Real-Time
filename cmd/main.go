package main

import (
	"log"
	"net/http"
	"os"

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

	// new service
	authSvc := service.NewAuthService(authRepo, userRepo)
	userSvc := service.NewUserService(userRepo)
	postSvc := service.NewPostService(postRepo)
	charSvc := service.NewChatService(chatRepo, userRepo)

	// new handler
	authHandler := handler.NewAuthHandler(authSvc, userSvc)
	postHandler := handler.NewPostHandler(postSvc)
	testHandler := handler.NewTestHandler(authSvc, userSvc)
	wsHandler := ws.NewHandleWs(charSvc, userSvc)
	chatHandler := handler.NewHandleChat(charSvc, userSvc)


	middleware := middleware.NewMiddleware(authSvc)
	router := handler.NewRouter(authHandler, postHandler, &testHandler, wsHandler, chatHandler, middleware)
	spa := spaHandler()
	http.ListenAndServe(":8081", spa(router))
	log.Println("Database Initialised")
}



func spaHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fileServer := http.FileServer(
			http.Dir("./web"),
		)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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