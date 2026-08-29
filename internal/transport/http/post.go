package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"realTime/internal/domain"
	"uuid"
)

type postService interface {
	CreatePost(context.Context, *domain.Post) error
	GetPost(context.Context, uuid.UUID, uuid.UUID) (*domain.PostInfo, error)
	GetPosts(context.Context, uuid.UUID, domain.PostFilter, int, uuid.UUID) ([]*domain.PostInfo, error)
}

type PostHandler struct {
	postSvc postService
}

func NewPostHandler(svc postService) *PostHandler {
	return &PostHandler{postSvc: svc}
}

func (p *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(uuid.UUID)
	request := domain.CreatePostRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w)
		return
	}

	post, err := domain.ValidatePostRequest(request)
	if err != nil {
		Error(err, w)
		return
	}
	post.AuthorID = userID
	err = p.postSvc.CreatePost(r.Context(), &post)
	if err != nil {
		Error(err, w)
		return
	};

	if err := json.NewEncoder(w).Encode(post.ID); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	};

}

func (p *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(uuid.UUID)
	query := r.URL.Query()

	
	cursor, err := domain.ValidatePostCursor(query.Get("cursor"))
	if err != nil {
		Error(err, w)
		return
	}

	
	filter, err := domain.ValidatePostFilter(query["category"], query.Get("liked"))
	if err != nil {
		Error(err, w)
		return
	}
	
	posts, err := p.postSvc.GetPosts(r.Context(), userID, filter, 20, cursor)
	if err != nil {
		Error(err, w)
		return
	}

	if posts == nil {
		posts = []*domain.PostInfo{};
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}

func (p *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(uuid.UUID)
	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		Error(domain.Error{Message: "invalid post id", Code: domain.BadFormatCode}, w)
		return
	}
	post, err := p.postSvc.GetPost(r.Context(), userID, postID)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}
