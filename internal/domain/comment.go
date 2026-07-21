package domain

import (
	"strings"

	"realTime/crypto"
)

type Comment struct {
	ID       crypto.UUID
	AuthorID crypto.UUID
	PostID   crypto.UUID
	ParentID crypto.UUID
	Content  string
}

type CommentRequest struct {
	PostID   string `json:"comment_id"`
	ParentID string `json:"parent_id"`
	Content  string `json:"content"`
}

func ValidateCommentRequest(commentRequest CommentRequest) (Comment, error) {
	comment := Comment{}
	commentRequest.Content = strings.TrimSpace(commentRequest.Content)
	if commentRequest.PostID == "" {
		return comment, Error{Message: "missing post id", Code: BadFormatCode}
	}

	if commentRequest.Content == "" {
		return comment, Error{Message: "missing comment content", Code: BadFormatCode}
	}

	if len(commentRequest.Content) < 10 || len(commentRequest.Content) > 4096 {
		return comment, Error{Message: "comment content length must be between 10 and 4096 chars", Code: BadFormatCode}
	}

	postID, err := crypto.ParseUUID(commentRequest.PostID)
	if err != nil {
		return comment, Error{Message: "invalid postID", Code: BadFormatCode}
	}
	if commentRequest.ParentID != "" {
		parentID, err := crypto.ParseUUID(commentRequest.ParentID)
		if err != nil {
			return comment, Error{Message: "invalid parent id", Code: BadFormatCode}
		}
		comment.ParentID = parentID
	}
	comment.PostID = postID
	comment.Content = commentRequest.Content
	return comment, nil
}
