package handler

import (
	"net/http"
)

type middleWare interface {
	Auth(http.Handler) http.Handler
}

func NewRouter(authHandler *AuthHandler, postHandler *PostHandler, commentHandler *CommentHandler, middleware middleWare) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/register", authHandler.Register)
	router.HandleFunc("POST /api/login", authHandler.Login)
	router.Handle("GET /api/posts", middleware.Auth(http.HandlerFunc(postHandler.GetPosts)))
	router.Handle("POST /api/posts", middleware.Auth(http.HandlerFunc(postHandler.CreatePost)))
	router.Handle("POST /api/comments", middleware.Auth(http.HandlerFunc(commentHandler.CreateComment)))
	router.Handle("PATCH /api/posts/{id}", middleware.Auth(http.HandlerFunc(postHandler.PatchPost)))
	router.Handle("GET /api/posts/{id}", middleware.Auth(http.HandlerFunc(postHandler.GetPost)))


	return router
}
