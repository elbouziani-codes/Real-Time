package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"realTime/crypto"
	"realTime/internal/domain"
)

type postService interface {
	CreatePost(context.Context, *domain.Post) error
	GetPost(context.Context, crypto.UUID) (*domain.Post, error)
	GetPosts(context.Context, int, int) ([]*domain.Post, error)
	PatchPost(context.Context, domain.PatchPostRequest, crypto.UUID) (error)
}

type PostHandler struct {
	postSvc postService
}

func NewPostHandler(svc postService) *PostHandler {
	return &PostHandler{postSvc: svc}
}

func (p *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)
	request := domain.CreatePostRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w)
		return
	}

	post, err := domain.ValueidatePostRequest(request)
	if err != nil {
		Error(err, w)
		return
	}
	post.AuthorID = userID
	err = p.postSvc.CreatePost(r.Context(), &post)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(post.ID); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}

func (p *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	
	/* request := domain.GetPostsRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w)
		return
	}

	err := domain.ValueidateGetPostsRequest(request)
	if err != nil {
		Error(err, w)
		return
	}*/
	// To  be dynamic later
	posts, err := p.postSvc.GetPosts(r.Context(), 20, 0)
	if err != nil {
		Error(err, w)
		return
	}
	
	
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}


func (p *PostHandler) PatchPost(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)
	postID, err := crypto.ParseUUID(r.PathValue("id")) 
	if err != nil {
		Error(domain.Error{Message: "invalid post id", Code: domain.BadFormatCode}, w)	
		return
	}

	request := domain.EditPostRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w)
		return
	}

	editPostObject, err := domain.ValueidateEditPostRequest(request)
	if err != nil {
		Error(err, w)
		return
	}

	post, err := p.postSvc.GetPost(r.Context(), postID)
	if err != nil {
		Error(err, w)
		return
	}
	
	if post.AuthorID != userID {
		Error(domain.Error{Message: "unauthorized action", Code: domain.UnauthorizedCode}, w)
		return
	}
	err = p.postSvc.PatchPost(r.Context(), editPostObject, postID)
	if err := json.NewEncoder(w).Encode(post.ID); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}
