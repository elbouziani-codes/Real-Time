package service

import (
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository"
)

type commentService struct {
	commentRepo CommentRepo
}

type CommentRepo interface {
	SaveComment(context.Context, domain.Comment) error
	GetComment(context.Context, crypto.UUID) (*domain.CommentInfo, error)
	GetComments(context.Context, int, int) ([]*domain.CommentInfo, error)	
	ExecPatchQuery(context.Context, string, crypto.UUID, []any) error 
}

func NewCommentService(repo CommentRepo) *commentService {
	return &commentService{commentRepo: repo}
}

func (p *commentService) CreateComment(ctx context.Context, comment *domain.Comment) error {
	id, err := crypto.GenerateUUID()
	if err != nil {
		return err
	}
	comment.ID = id
	err = p.commentRepo.SaveComment(ctx, *comment)
	if err != nil {
		return err
	}
	return nil
}

func (p *commentService) GetComment(ctx context.Context, commentID crypto.UUID) (*domain.CommentInfo, error) {
	return  p.commentRepo.GetComment(ctx, commentID) 
}


func (p *commentService) PatchComment(ctx context.Context, editObject domain.PatchCommentRequest, commentID crypto.UUID) error {
	query := repository.ConstructPatchQuery("comments", "id",  editObject.FilledKeys)			
	return p.commentRepo.ExecPatchQuery(ctx, query, commentID, editObject.FilledValues) 
}


func (p *commentService) GetComments(ctx context.Context, limit, offset int) ([]*domain.CommentInfo, error) {
	return p.commentRepo.GetComments(context.Background(), limit, offset)
}


