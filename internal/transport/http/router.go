package handler

import (
	"net/http"
)

type middleWare interface {
	Auth(http.Handler) http.Handler
}

func NewRouter(authHandler *AuthHandler, postHandler *PostHandler, commentHandler *CommentHandler, reactionHandler *ReactionHandler, middleware middleWare) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/register", authHandler.Register)
	router.HandleFunc("POST /api/login", authHandler.Login)
	router.Handle("GET /api/posts", middleware.Auth(http.HandlerFunc(postHandler.GetPosts)))
	router.Handle("POST /api/posts", middleware.Auth(http.HandlerFunc(postHandler.CreatePost)))
	router.Handle("POST /api/comments", middleware.Auth(http.HandlerFunc(commentHandler.CreateComment)))
	router.Handle("DELETE /api/comments/{id}", middleware.Auth(http.HandlerFunc(commentHandler.DeleteComment)))
	router.Handle("PATCH /api/posts/{id}", middleware.Auth(http.HandlerFunc(postHandler.PatchPost)))
	router.Handle("PATCH /api/comments/{id}", middleware.Auth(http.HandlerFunc(commentHandler.PatchComment)))
	router.Handle("GET /api/posts/{id}", middleware.Auth(http.HandlerFunc(postHandler.GetPost)))
	router.Handle("GET /api/me", middleware.Auth(http.HandlerFunc(authHandler.Me)))

	router.Handle("GET /api/comments/{id}", middleware.Auth(http.HandlerFunc(commentHandler.GetComments)))
	router.Handle("GET /api/reactions/{parent_id}", middleware.Auth(http.HandlerFunc(reactionHandler.GetReactions)))
	router.Handle("POST /api/reactions", middleware.Auth(http.HandlerFunc(reactionHandler.React)))
	router.Handle("PATCH /api/reactions/{id}", middleware.Auth(http.HandlerFunc(reactionHandler.UpdateReaction)))




	return router
}
