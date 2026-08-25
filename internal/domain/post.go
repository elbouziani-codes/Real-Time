package domain

import (
	"fmt"
	"strings"
	"uuid"
)

const MaxPostCategories = 5

type Post struct {
	ID         uuid.UUID
	AuthorID   uuid.UUID
	Categories []uuid.UUID
	Title      string
	Content    string
	CreatedAt  int
}

type PostInfo struct {
	ID         uuid.UUID
	Author     UserProfile
	Title      string
	Content    string
	Categories []Category
	Likes      int
	DisLikes   int
	LikeInfo   LikeInfo
	CreatedAt  int
}

type LikeInfo struct {
	ID     uuid.UUID
	IsLike bool
}

type CreatePostRequest struct {
	Title      string   `json:"title"`
	Categories []string `json:"category_ids"`
	Content    string   `json:"content"`
}

func ValidatePostRequest(request CreatePostRequest) (Post, error) {
	post := Post{}
	request.Title = strings.TrimSpace(request.Title)
	request.Content = strings.TrimSpace(request.Content)

	if request.Title == "" {
		return post, Error{Message: "missing post title", Code: BadFormatCode}
	}
	if request.Content == "" {
		return post, Error{Message: "missing post content", Code: BadFormatCode}
	}

	if len(request.Title) < 1 || len(request.Title) > 100 {
		return post, Error{Message: "post title length must be between 1 and 100 chars", Code: BadFormatCode}
	}
	if len(request.Content) < 1 || len(request.Content) > 4096 {
		return post, Error{Message: "post content length must be between 1 and 4096 chars", Code: BadFormatCode}
	}
	if len(request.Categories) == 0 {
		return post, Error{Message: "missing post category", Code: BadFormatCode}
	}
	if len(request.Categories) > MaxPostCategories {
		return post, Error{Message: fmt.Sprintf("a post accepts at most %d categories", MaxPostCategories), Code: BadFormatCode}
	}

	
	seen := make(map[uuid.UUID]struct{}, len(request.Categories))
	for _, rawID := range request.Categories {
		categoryID, err := uuid.Parse(rawID)
		if err != nil {
			return post, Error{Message: "invalid post category", Code: BadFormatCode}
		}
		if _, duplicated := seen[categoryID]; duplicated {
			return post, Error{Message: "duplicated post category", Code: BadFormatCode}
		}
		seen[categoryID] = struct{}{}
		post.Categories = append(post.Categories, categoryID)
	}

	post.Title = request.Title
	post.Content = request.Content
	return post, nil
}



type PostFilter struct {
	Categories []uuid.UUID
	LikedOnly  bool
}


func ValidatePostFilter(rawCategories []string, rawLiked string) (PostFilter, error) {
	filter := PostFilter{}

	switch strings.TrimSpace(rawLiked) {
	case "", "false":
	case "true":
		filter.LikedOnly = true
	default:
		return filter, Error{Message: "liked must be true or false", Code: BadFormatCode}
	}


	
	seen := make(map[uuid.UUID]struct{}, len(rawCategories))
	for _, rawID := range rawCategories {
		categoryID, err := uuid.Parse(rawID)
		if err != nil {
			return PostFilter{}, Error{Message: "invalid category filter", Code: BadFormatCode}
		}
		if _, duplicated := seen[categoryID]; duplicated {
			continue
		}
		seen[categoryID] = struct{}{}
		filter.Categories = append(filter.Categories, categoryID)
	}

	return filter, nil
}


func ValidatePostCursor(rawCursor string) (uuid.UUID, error) {
	rawCursor = strings.TrimSpace(rawCursor)
	if rawCursor == "" {
		return uuid.Nil(), nil
	}

	cursor, err := uuid.Parse(rawCursor)
	if err != nil {
		return uuid.Nil(), Error{Message: "invalid cursor", Code: BadFormatCode}
	}
	return cursor, nil
}
