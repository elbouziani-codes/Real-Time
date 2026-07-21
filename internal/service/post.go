package service

import (
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository"
)

type postService struct {
	postRepo PostRepo
}

type PostRepo interface {
	SavePost(context.Context, domain.Post) error
	GetPost(context.Context, crypto.UUID) (*domain.Post, error)
	GetPosts(context.Context, int, int) ([]*domain.Post, error)	
	ExecPatchQuery(context.Context, string, crypto.UUID, []any) error 
}

func NewPostService(repo PostRepo) *postService {
	return &postService{postRepo: repo}
}

func (p *postService) CreatePost(ctx context.Context, post *domain.Post) error {
	id, err := crypto.GenerateUUID()
	if err != nil {
		return err
	}
	post.ID = id
	err = p.postRepo.SavePost(ctx, *post)
	if err != nil {
		return err
	}
	return nil
}

func (p *postService) GetPost(ctx context.Context, postID crypto.UUID) (*domain.Post, error) {
	return  p.postRepo.GetPost(ctx, postID) 
}


func (p *postService) PatchPost(ctx context.Context, editObject domain.PatchPostRequest, postID crypto.UUID) error {
	query := repository.ConstructPatchQuery("posts", "id",  editObject.FilledKeys)			
	return p.postRepo.ExecPatchQuery(ctx, query, postID, editObject.FilledValues) 
}


func (p *postService) GetPosts(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	return p.postRepo.GetPosts(context.Background(), limit, offset)
}


