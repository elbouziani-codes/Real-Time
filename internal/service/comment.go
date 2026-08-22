package service

import (
	"context"
	"realTime/internal/domain"
	"uuid"
)

type commentService struct {
	commentRepo CommentRepo
}

type CommentRepo interface {
	SaveComment(context.Context, domain.Comment) error
	GetComment(context.Context, uuid.UUID, uuid.UUID) (*domain.CommentInfo, error)
	DeleteComment(context.Context, uuid.UUID) error
	GetComments(context.Context, uuid.UUID, uuid.UUID) ([]*domain.CommentInfo, error)
}

func NewCommentService(repo CommentRepo) *commentService {
	return &commentService{commentRepo: repo}
}

func (p *commentService) CreateComment(ctx context.Context, comment *domain.Comment) error {
	comment.ID = uuid.NewV4()
	err := p.commentRepo.SaveComment(ctx, *comment)
	if err != nil {
		return err
	}
	return nil
}

func (p *commentService) GetComment(ctx context.Context, userID, commentID uuid.UUID) (*domain.CommentInfo, error) {
	return p.commentRepo.GetComment(ctx, userID, commentID)
}

func (p *commentService) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	return p.commentRepo.DeleteComment(ctx, commentID)
}

func (p *commentService) GetComments(ctx context.Context, userID, parentID uuid.UUID) ([]*domain.CommentInfo, error) {
	return p.commentRepo.GetComments(ctx, userID, parentID)
}
