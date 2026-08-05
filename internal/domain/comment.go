package domain

import (
	"fmt"
	"realTime/crypto"
	"strings"
)

type Comment struct {
	ID        crypto.UUID
	ParentID  crypto.UUID
	AuthorID  crypto.UUID
	Content   string
}


type CommentInfo struct {
	ID        crypto.UUID
	ParentID  crypto.UUID
	Author 	  UserProfile
	Likes 	  int
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
	fmt.Println(request.ParentID)
	if err != nil {
		return comment, Error{Message: "invalid parent_id", Code: BadFormatCode}
	}

	comment.Content = request.Content
	return comment, nil
}


type EditCommentRequest struct {	
	Content string `json:"content"`
}

type PatchCommentRequest struct {
	FilledKeys []string
	FilledValues []any
}


func NewEditCommentRequest(request EditCommentRequest) (PatchCommentRequest, error) {
	comment := PatchCommentRequest{}
	request.Content = strings.TrimSpace(request.Content)

		
	if request.Content != "" {
		if ( len(request.Content) < 10 || len(request.Content) > 4096 && request.Content != "") {
			return comment, Error{Message: "comment content length must be between 10 and 4096 chars", Code: BadFormatCode}
		}		
		comment.FilledKeys = append(comment.FilledKeys, "content") 
		comment.FilledValues = append(comment.FilledValues, request.Content) 

	} else {
		return comment, Error{Message: "comment content length must be between 10 and 4096 chars", Code: BadFormatCode}
	}	
		
		
	return comment, nil
}



type GetCommentsRequest struct {
	Offset  int `json:"offset"`
	Limit   int  `json:"limit"`
}

func ValueidateGetCommentsRequest(request GetCommentsRequest) (error) {
	if request.Offset <= 0 || request.Limit <= 0 {
			return Error{Message: "invalid filters", Code: BadFormatCode}		
	}	

	if  request.Limit >= 50 {
			return Error{Message: "nah not that time", Code: BadFormatCode}		
	}
	return nil  
}


