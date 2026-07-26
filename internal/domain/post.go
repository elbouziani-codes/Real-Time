package domain

import (
	"realTime/crypto"
	"strings"
)

type Post struct {
	ID       crypto.UUID
	AuthorID crypto.UUID
	Title    string
	Content  string
	CreatedAt int
	UpdatedAt int
}


type PostInfo struct {
	ID        crypto.UUID
	Author UserProfile
	Title     string
	Content   string
	Likes 	  int
	DisLikes   int
	CreatedAt int
	UpdatedAt int
}




type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func ValueidatePostRequest(request CreatePostRequest) (Post, error) {
	post := Post{}
	request.Title = strings.TrimSpace(request.Title)
	request.Content = strings.TrimSpace(request.Content)

	if request.Title == "" {
		return post, Error{Message: "missing post title", Code: BadFormatCode}
	}
	if request.Content == "" {
		return post, Error{Message: "missing post content", Code: BadFormatCode}
	}

	if len(request.Title) < 10 || len(request.Title) > 100 {
		return post, Error{Message: "post title length must be between 10 and 100 chars", Code: BadFormatCode}
	}
	if len(request.Content) < 10 || len(request.Content) > 4096 {
		return post, Error{Message: "post content length must be between 10 and 4096 chars", Code: BadFormatCode}
	}

	post.Title = request.Title
	post.Content = request.Content
	return post, nil
}


type EditPostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PatchPostRequest struct {
	FilledKeys []string
	FilledValues []any
}


func ValueidateEditPostRequest(request EditPostRequest) (PatchPostRequest, error) {
	post := PatchPostRequest{}
	request.Title = strings.TrimSpace(request.Title)
	request.Content = strings.TrimSpace(request.Content)

	
	if request.Title != "" {
		if (len(request.Title) < 10 || len(request.Title) > 100) {
			return post, Error{Message: "post title length must be between 10 and 100 chars", Code: BadFormatCode}
		}			
		post.FilledKeys = append(post.FilledKeys, "title") 
		post.FilledValues= append(post.FilledValues, request.Title) 
	}	
	if request.Content != "" {
		if ( len(request.Content) < 10 || len(request.Content) > 4096 && request.Content != "") {
			return post, Error{Message: "post content length must be between 10 and 4096 chars", Code: BadFormatCode}
		}		
		post.FilledKeys = append(post.FilledKeys, "content") 
		post.FilledValues = append(post.FilledValues, request.Content) 

	}	
		
		
	return post, nil
}



type GetPostsRequest struct {
	Offset  int `json:"offset"`
	Limit   int  `json:"limit"`
}

func ValueidateGetPostsRequest(request GetPostsRequest) (error) {
	if request.Offset <= 0 || request.Limit <= 0 {
			return Error{Message: "invalid filters", Code: BadFormatCode}		
	}	

	if  request.Limit >= 50 {
			return Error{Message: "nah not that time", Code: BadFormatCode}		
	}
	return nil  
}


