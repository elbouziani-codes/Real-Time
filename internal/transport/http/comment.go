package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"realTime/internal/domain"
	"uuid"
)

type commentService interface {
	CreateComment(context.Context, *domain.Comment) error
	// GetComment(context.Context, uuid.UUID, uuid.UUID) (*domain.CommentInfo, error)
	// DeleteComment(context.Context, uuid.UUID) error
	GetComments(context.Context, uuid.UUID, uuid.UUID) ([]*domain.CommentInfo, error)
}

type CommentHandler struct {
	commentSvc commentService
	postSvc    postService
}

func NewCommentHandler(commentSvc commentService, postSvc postService) *CommentHandler {
	return &CommentHandler{commentSvc: commentSvc, postSvc: postSvc}
}

func (c *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(uuid.UUID)
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

// func (c *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
// 	userIDAny := r.Context().Value("user_id")
// 	userID := userIDAny.(uuid.UUID)
//
// 	commentID, err := uuid.Parse(r.PathValue("id"))
// 	if err != nil {
// 		Error(domain.Error{Message: "invalid comment id", Code: domain.BadFormatCode}, w)
// 		return
// 	}
// 	comment, err := c.commentSvc.GetComment(r.Context(), userID, commentID)
// 	if err != nil {
// 		Error(err, w)
// 		return
// 	}
// 	if comment.Author.ID.String() != userID.String() {
// 		Error(domain.Error{Message: "unauthorized action", Code: domain.UnauthorizedCode}, w)
// 		return
// 	}
// 	err = c.commentSvc.DeleteComment(r.Context(), commentID)
// 	if err != nil {
// 		Error(err, w)
// 		return
// 	}
//
// 	if err := json.NewEncoder(w).Encode("success"); err != nil {
// 		http.Error(w, "uknown error", http.StatusInternalServerError)
// 		return
// 	}
//
// }

func (c *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(uuid.UUID)

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		Error(domain.Error{Message: "invalid id", Code: domain.BadFormatCode}, w)
		return
	}

	comments, err := c.commentSvc.GetComments(r.Context(), userID, id)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}
