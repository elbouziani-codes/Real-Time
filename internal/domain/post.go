package domain

import (
	"fmt"
	"realTime/crypto"
	"strings"
)

const MaxPostCategories = 5

type Post struct {
	ID         crypto.UUID
	AuthorID   crypto.UUID
	Categories []crypto.UUID
	Title      string
	Content    string
	CreatedAt  int
	UpdatedAt  int
}

type PostInfo struct {
	ID         crypto.UUID
	Author     UserProfile
	Title      string
	Content    string
	Categories []Category
	Likes      int
	DisLikes   int
	LikeInfo   LikeInfo
	CreatedAt  int
	UpdatedAt  int
}

type LikeInfo struct {
	ID     crypto.UUID
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

	if len(request.Title) < 10 || len(request.Title) > 100 {
		return post, Error{Message: "post title length must be between 10 and 100 chars", Code: BadFormatCode}
	}
	if len(request.Content) < 10 || len(request.Content) > 4096 {
		return post, Error{Message: "post content length must be between 10 and 4096 chars", Code: BadFormatCode}
	}
	if len(request.Categories) == 0 {
		return post, Error{Message: "missing post category", Code: BadFormatCode}
	}
	if len(request.Categories) > MaxPostCategories {
		return post, Error{Message: fmt.Sprintf("a post accepts at most %d categories", MaxPostCategories), Code: BadFormatCode}
	}

	// Reject duplicates here so a repeated id is a 400 rather than a primary
	// key violation surfacing from post_categories.
	seen := make(map[crypto.UUID]struct{}, len(request.Categories))
	for _, rawID := range request.Categories {
		categoryID, err := crypto.ParseUUID(rawID)
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

// MaxFilterCategories bounds the IN clause the category filter expands into, so
// a caller cannot force an arbitrarily large query by repeating the parameter.
const MaxFilterCategories = 20

// PostFilter narrows a post listing. A zero value means "no filtering", so the
// unfiltered listing keeps working unchanged.
type PostFilter struct {
	Categories []crypto.UUID
	LikedOnly  bool
}

// ValidatePostFilter turns the raw query parameters into a PostFilter. Unknown
// category ids are not an error: they simply match no post, which keeps a
// stale bookmark from turning into a 404.
func ValidatePostFilter(rawCategories []string, rawLiked string) (PostFilter, error) {
	filter := PostFilter{}

	switch strings.TrimSpace(rawLiked) {
	case "", "false":
	case "true":
		filter.LikedOnly = true
	default:
		return filter, Error{Message: "liked must be true or false", Code: BadFormatCode}
	}

	if len(rawCategories) > MaxFilterCategories {
		return filter, Error{Message: fmt.Sprintf("a listing accepts at most %d category filters", MaxFilterCategories), Code: BadFormatCode}
	}

	// Duplicates would only repeat a placeholder without changing the result,
	// so drop them rather than reject the request.
	seen := make(map[crypto.UUID]struct{}, len(rawCategories))
	for _, rawID := range rawCategories {
		categoryID, err := crypto.ParseUUID(rawID)
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

// ValidatePostCursor parses the id of the last post the client already holds.
// Paging on that post's sort position rather than a row count keeps a page from
// skipping or repeating posts when others are published mid-scroll. An empty
// value asks for the first page, so an opening request carries no cursor.
func ValidatePostCursor(rawCursor string) (crypto.UUID, error) {
	rawCursor = strings.TrimSpace(rawCursor)
	if rawCursor == "" {
		return crypto.Nil, nil
	}

	cursor, err := crypto.ParseUUID(rawCursor)
	if err != nil {
		return crypto.Nil, Error{Message: "invalid cursor", Code: BadFormatCode}
	}
	return cursor, nil
}
