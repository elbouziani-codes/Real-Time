package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"realTime/internal/domain"
	"uuid"
)

type reactionService interface {
	CreateReaction(context.Context, *domain.Reaction) error
	GetReaction(context.Context, uuid.UUID) (*domain.ReactionInfo, error)
	GetReactions(context.Context, uuid.UUID) ([]*domain.ReactionInfo, error)
	PatchReaction(context.Context, *domain.ReactionInfo, bool) error
}

type ReactionHandler struct {
	reactionSvc reactionService
	postSvc     postService
}

func NewReactionHandler(reactionSvc reactionService, postSvc postService) *ReactionHandler {
	return &ReactionHandler{reactionSvc: reactionSvc, postSvc: postSvc}
}

func (c *ReactionHandler) React(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(uuid.UUID)
	request := domain.ReactionRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w)
		return
	}

	reaction, err := domain.NewReactionRequest(request)
	if err != nil {
		Error(err, w)
		return
	}
	reaction.AuthorID = userID
	err = c.reactionSvc.CreateReaction(r.Context(), &reaction)
	if err != nil {
		Error(err, w)
		return
	}
	if err := json.NewEncoder(w).Encode(reaction.ID); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}

func (c *ReactionHandler) GetReactions(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("parent_id"))
	if err != nil {
		Error(domain.Error{Message: "invalid id", Code: domain.BadFormatCode}, w)
		return
	}

	reactions, err := c.reactionSvc.GetReactions(r.Context(), id)
	if err != nil {
		Error(err, w)
		return
	}

	if err := json.NewEncoder(w).Encode(reactions); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}

}

func (c *ReactionHandler) UpdateReaction(w http.ResponseWriter, r *http.Request) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(uuid.UUID)
	reactionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		Error(domain.Error{Message: "invalid reaction id", Code: domain.BadFormatCode}, w)
		return
	}

	request := domain.UpdateReactionRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(domain.Error{Message: "invalid json", Code: domain.BadFormatCode}, w)
		return
	}

	react, err := c.reactionSvc.GetReaction(r.Context(), reactionID)
	if err != nil {
		Error(err, w)
		return
	}

	if react.Author.ID.String() != userID.String() {
		Error(domain.Error{Message: "unauthorized action", Code: domain.UnauthorizedCode}, w)
		return
	}
	err = c.reactionSvc.PatchReaction(r.Context(), react, request.IsLike)
	if err != nil {
		Error(err, w)
		return
	}
	if err := json.NewEncoder(w).Encode("success"); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}
}
