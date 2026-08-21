package domain

import (
	"realTime/crypto"
	"strings"
)

type Comment struct {
	ID       crypto.UUID
	ParentID crypto.UUID
	AuthorID crypto.UUID
	Content  string
}

type CommentInfo struct {
	ID        crypto.UUID
	ParentID  crypto.UUID
	Author    UserProfile
	Likes     int
	DisLike   int
	Content   string
	LikeInfo  LikeInfo
	CreatedAt int
	UpdatedAt int
}

type CreateCommentRequest struct {
	ParentID string `json:"parent_id"`
	Content  string `json:"content"`
}

func NewCommentRequest(request CreateCommentRequest) (Comment, error) {
	comment := Comment{}
	request.Content = strings.TrimSpace(request.Content)

	if request.ParentID == "" {
		return comment, Error{Message: "missing comment.Parent_id", Code: BadFormatCode}
	}
	if request.Content == "" {
		return comment, Error{Message: "missing comment content", Code: BadFormatCode}
	}

	if len(request.Content) < 1 || len(request.Content) > 4096 {
		return comment, Error{Message: "comment content length must be between 1 and 4096 chars", Code: BadFormatCode}
	}
	var err error
	comment.ParentID, err = crypto.ParseUUID(request.ParentID)
	if err != nil {
		return comment, Error{Message: "invalid parent_id", Code: BadFormatCode}
	}

	comment.Content = request.Content
	return comment, nil
}
