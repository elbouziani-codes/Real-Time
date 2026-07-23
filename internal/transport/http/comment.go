package handler

import (
	"strconv"
	"context"
	"encoding/json"
	"net/http"
	"realTime/crypto"
	"realTime/internal/domain"
)

type commentService interface {
	CreateComment(context.Context, *domain.Comment) error
	GetComment(context.Context, crypto.UUID) (*domain.CommentInfo, error)
	GetComments(context.Context, int, int) ([]*domain.CommentInfo, error)
	PatchComment(context.Context, domain.PatchCommentRequest, crypto.UUID) (error)
}

type CommentHandler struct {
	commentSvc commentService
	postSvc postService
}

func NewCommentHandler(commentSvc commentService, postSvc postService) *CommentHandler {
		return &CommentHandler{commentSvc: commentSvc, postSvc: postSvc}
}

func (p *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)
	request := domain.CreateCommentRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w)
		return
	}

	

	comment, err := domain.NewCommentRequest(request)
	if err != nil {
		Error(err, w)
		return
	}
	comment.AuthorID = userID
	err = p.commentSvc.CreateComment(r.Context(), &comment)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(comment.ID); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}

func (p *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	
	query := r.URL.Query() 
	n := query.Get("offset") 
	offset, _ := strconv.Atoi(n) // I dont need to check error because if it failed then it will be 0   			
	if offset < 0 {
		offset = 0 //fallbacking to 0
	}


	comments, err := p.commentSvc.GetComments(r.Context(), 20, offset)
	if err != nil {
		Error(err, w)
		return
	}
	
	
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}

func (p *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := crypto.ParseUUID(r.PathValue("id")) 
	if err != nil {
		Error(domain.Error{Message: "invalid comment id", Code: domain.BadFormatCode}, w)	
		return
	}
	comment, err := p.commentSvc.GetComment(r.Context(), commentID)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}
	
}

func (p *CommentHandler) PatchComment(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)
	commentID, err := crypto.ParseUUID(r.PathValue("id")) 
	if err != nil {
		Error(domain.Error{Message: "invalid comment id", Code: domain.BadFormatCode}, w)	
		return
	}

	request := domain.EditCommentRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w)
		return
	}

	editCommentObject, err := domain.NewEditCommentRequest(request)
	if err != nil {
		Error(err, w)
		return
	}

	comment, err := p.commentSvc.GetComment(r.Context(), commentID)
	if err != nil {
		Error(err, w)
		return
	}
	
	if comment.Author.ID.Value != userID.Value {
		Error(domain.Error{Message: "unauthorized action", Code: domain.UnauthorizedCode}, w)
		return
	}
	err = p.commentSvc.PatchComment(r.Context(), editCommentObject, commentID)
	if err := json.NewEncoder(w).Encode(comment.ID); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}
