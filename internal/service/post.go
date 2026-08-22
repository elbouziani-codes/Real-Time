package service

import (
	"context"
	"realTime/internal/domain"
	"uuid"
)

type postService struct {
	postRepo PostRepo
}

type PostRepo interface {
	SavePost(context.Context, domain.Post) error
	GetPost(context.Context, uuid.UUID, uuid.UUID) (*domain.PostInfo, error)
	GetPosts(context.Context, uuid.UUID, domain.PostFilter, int, uuid.UUID) ([]*domain.PostInfo, error)
}

func NewPostService(repo PostRepo) *postService {
	return &postService{postRepo: repo}
}

func (p *postService) CreatePost(ctx context.Context, post *domain.Post) error {
	post.ID = uuid.NewV4()
	err := p.postRepo.SavePost(ctx, *post)
	if err != nil {
		return err
	}
	return nil
}

func (p *postService) GetPost(ctx context.Context, userID, postID uuid.UUID) (*domain.PostInfo, error) {
	return p.postRepo.GetPost(ctx, userID, postID)
}

func (p *postService) GetPosts(ctx context.Context, userID uuid.UUID, filter domain.PostFilter, limit int, cursor uuid.UUID) ([]*domain.PostInfo, error) {
	return p.postRepo.GetPosts(ctx, userID, filter, limit, cursor)
}
