package handler

import (
	"net/http"

	ws "realTime/internal/transport/webSocket"
)

type middleWare interface {
	Auth(http.Handler) http.Handler
}

func NewRouter(authHandler *AuthHandler, postHandler *PostHandler, testHandler *TestHandler, wsHander *ws.HandlerWs, chatHandler *HandlerChat, middleware middleWare) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/register", authHandler.Register)
	router.HandleFunc("POST /api/login", authHandler.Login)
	router.Handle("POST /api/getMessage", middleware.Auth(http.HandlerFunc(chatHandler.Chat)))
	router.Handle("GET /api/posts", middleware.Auth(http.HandlerFunc(postHandler.GetPosts)))
	router.Handle("GET /api/ws", middleware.Auth(http.HandlerFunc(wsHander.ChatWs)))
	router.Handle("POST /api/posts", middleware.Auth(http.HandlerFunc(postHandler.CreatePost)))
	router.Handle("PATCH /api/posts/{id}", middleware.Auth(http.HandlerFunc(postHandler.PatchPost)))
	
	router.Handle("/api/test", middleware.Auth(http.HandlerFunc(testHandler.Test)))
	fs := http.FileServer(
		http.Dir("./web"),
	)

	http.Handle("/", fs)
	
	return router
}
