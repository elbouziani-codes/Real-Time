package service

import (
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
)

type reactionService struct {
	reactionRepo ReactionRepo
}

type ReactionRepo interface {
	SaveReaction(context.Context, domain.Reaction) error
	GetReaction(context.Context, crypto.UUID) (*domain.ReactionInfo, error)
	GetReactionByUserAndParent(context.Context, crypto.UUID, crypto.UUID) (*domain.ReactionInfo, error)
	DeleteReaction(context.Context, crypto.UUID) (error)
	GetReactions(context.Context, crypto.UUID) ([]*domain.ReactionInfo, error)	
	UpdateReaction(context.Context, crypto.UUID, bool) error 
}

func NewReactionService(repo ReactionRepo) *reactionService {
	return &reactionService{reactionRepo: repo}
}

func (p *reactionService) CreateReaction(ctx context.Context, reaction *domain.Reaction) error {
	// first check if comment exist then I will return conflict 
	_, err := p.reactionRepo.GetReactionByUserAndParent(ctx, reaction.ParentID, reaction.AuthorID) 
	if err == nil {
			return domain.Error{Message:"You already raected", Code: domain.ConflictCode}	
	}
	id, err := crypto.GenerateUUID()
	if err != nil {
		return err
	}
	reaction.ID = id
	err = p.reactionRepo.SaveReaction(ctx, *reaction)
	if err != nil {
		return err
	}
	return nil
}


func (p *reactionService) DeleteReaction(ctx context.Context, reactionID crypto.UUID) (error) {
	return  p.reactionRepo.DeleteReaction(ctx, reactionID) 
}


func (p *reactionService) PatchReaction(ctx context.Context, react *domain.ReactionInfo, isLike bool) error {
	
	if react.IsLike == isLike {
		return p.DeleteReaction(ctx, react.ID) 
	}
	return p.reactionRepo.UpdateReaction(ctx, react.ID, isLike) 
}

func (p *reactionService) GetReaction(ctx context.Context, reactionID crypto.UUID) (*domain.ReactionInfo, error) {
	return p.reactionRepo.GetReaction(context.Background(), reactionID)
}

func (p *reactionService) GetReactions(ctx context.Context, parentID crypto.UUID) ([]*domain.ReactionInfo, error) {
	return p.reactionRepo.GetReactions(context.Background(), parentID)
}
