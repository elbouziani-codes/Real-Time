package service

import (
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
)

type categoryService struct {
	categoryRepo CategoryRepo
}

type CategoryRepo interface {
	GetCategories(context.Context) ([]*domain.Category, error)
	GetCategory(context.Context, crypto.UUID) (*domain.Category, error)
}

func NewCategoryService(repo CategoryRepo) *categoryService {
	return &categoryService{categoryRepo: repo}
}

func (c *categoryService) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	return c.categoryRepo.GetCategories(ctx)
}

func (c *categoryService) GetCategory(ctx context.Context, categoryID crypto.UUID) (*domain.Category, error) {
	return c.categoryRepo.GetCategory(ctx, categoryID)
}
