package domain

import (
	"realTime/crypto"
)

type Reaction struct {
	ID       crypto.UUID
	ParentID crypto.UUID
	AuthorID crypto.UUID
	IsLike   bool
}

type ReactionInfo struct {
	ID        crypto.UUID
	ParentID  crypto.UUID
	Author    UserProfile
	IsLike    bool
	CreatedAt int
	UpdatedAt int
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
	react.ParentID, err = crypto.ParseUUID(request.ParentID)
	if err != nil {
		return react, Error{Message: "invalid parent_id", Code: BadFormatCode}
	}

	react.IsLike = request.IsLike
	return react, nil
}

type UpdateReactionRequest struct {
	IsLike bool `json:"is_like"`
}
