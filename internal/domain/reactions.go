package domain

import (
	"uuid"
)

type Reaction struct {
	ID       uuid.UUID
	ParentID uuid.UUID
	AuthorID uuid.UUID
	IsLike   bool
}

type ReactionInfo struct {
	ID        uuid.UUID
	ParentID  uuid.UUID
	Author    UserProfile
	IsLike    bool
	CreatedAt int
}

type ReactionRequest struct {
	ParentID string `json:"parent_id"`
	IsLike   bool   `json:"is_like"` // just fallback to false if I ddi use pointer .. I will not handle it
}

func NewReactionRequest(request ReactionRequest) (Reaction, error) {
	react := Reaction{}

	if request.ParentID == "" {
		return react, Error{Message: "missing react parent_id", Code: BadFormatCode}
	}

	var err error
	react.ParentID, err = uuid.Parse(request.ParentID)
	if err != nil {
		return react, Error{Message: "invalid parent_id", Code: BadFormatCode}
	}

	react.IsLike = request.IsLike
	return react, nil
}

type UpdateReactionRequest struct {
	IsLike bool `json:"is_like"`
}
