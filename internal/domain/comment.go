package domain

import (
	"strings"
	"uuid"
)

type Comment struct {
	ID       uuid.UUID
	PostID   uuid.UUID
	AuthorID uuid.UUID
	Content  string
}

type CommentInfo struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	Author    UserProfile
	Content   string
	CreatedAt int
}

type CreateCommentRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

func NewCommentRequest(request CreateCommentRequest) (Comment, error) {
	comment := Comment{}
	request.Content = strings.TrimSpace(request.Content)

	if request.PostID == "" {
		return comment, Error{Message: "missing comment.post_id", Code: BadFormatCode}
	}
	if request.Content == "" {
		return comment, Error{Message: "missing comment content", Code: BadFormatCode}
	}

	if len(request.Content) < 1 || len(request.Content) > 4096 {
		return comment, Error{Message: "comment content length must be between 1 and 4096 chars", Code: BadFormatCode}
	}
	var err error
	comment.PostID, err = uuid.Parse(request.PostID)
	if err != nil {
		return comment, Error{Message: "invalid post_id", Code: BadFormatCode}
	}

	comment.Content = request.Content
	return comment, nil
}
