package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"realTime/crypto"
	"realTime/internal/domain"
)

type commentService interface {
	CreateComment(context.Context, *domain.Comment) error
	GetComment(context.Context, crypto.UUID, crypto.UUID) (*domain.CommentInfo, error)
	DeleteComment(context.Context, crypto.UUID) (error)
	GetComments(context.Context, crypto.UUID, crypto.UUID, int) ([]*domain.CommentInfo, error)
	PatchComment(context.Context, domain.PatchCommentRequest, crypto.UUID) (error)
}

type CommentHandler struct {
	commentSvc commentService
	postSvc postService
}

func NewCommentHandler(commentSvc commentService, postSvc postService) *CommentHandler {
		return &CommentHandler{commentSvc: commentSvc, postSvc: postSvc}
}

func (c *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
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
	err = c.commentSvc.CreateComment(r.Context(), &comment)
	if err != nil {
		Error(err, w)
		return
	}
	if err := json.NewEncoder(w).Encode(comment.ID); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}


func (c *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)
	

	commentID, err := crypto.ParseUUID(r.PathValue("id")) 
	if err != nil {
		Error(domain.Error{Message: "invalid comment id", Code: domain.BadFormatCode}, w)	
		return
	}
	comment, err := c.commentSvc.GetComment(r.Context(), userID, commentID)
	if err != nil {
		Error(err, w)
		return
	}
	if comment.Author.ID.Value != userID.Value {
			Error(domain.Error{Message: "unauthorized action", Code: domain.UnauthorizedCode}, w)	
		return
	}
	err =c.commentSvc.DeleteComment(r.Context(), commentID)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode("success"); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}

func (c *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)

	id, err := crypto.ParseUUID(r.PathValue("id")) 
	cursor := r.PathValue("cursor") 
	
	if err != nil {
			Error(domain.Error{Message: "invalid id", Code: domain.BadFormatCode}, w)
			return
	}		
	
	cur, err := strconv.Atoi(cursor)		
	if err != nil {
		Error(domain.Error{Message: "invalid cursor", Code: domain.BadFormatCode}, w)
		return
	
	}
	comments, err := c.commentSvc.GetComments(r.Context(), userID, id, cur)
	if err != nil {
		Error(err, w)
		return
	}
	
	
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}
}

func (c *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(crypto.UUID)

	commentID, err := crypto.ParseUUID(r.PathValue("id")) 
	if err != nil {
		Error(domain.Error{Message: "invalid comment id", Code: domain.BadFormatCode}, w)	
		return
	}
	comment, err :=c.commentSvc.GetComment(r.Context(), userID, commentID)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}
	
}

func (c *CommentHandler) PatchComment(w http.ResponseWriter, r *http.Request) {
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

	comment, err :=c.commentSvc.GetComment(r.Context(), userID, commentID)
	if err != nil {
		Error(err, w)
		return
	}
	
	if comment.Author.ID.Value != userID.Value {
		Error(domain.Error{Message: "unauthorized action", Code: domain.UnauthorizedCode}, w)
		return
	}
	err =c.commentSvc.PatchComment(r.Context(), editCommentObject, commentID)
	if err := json.NewEncoder(w).Encode(comment.ID); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}
