package service

import (
	"context"
	"realTime/internal/domain"
)

type categoryService struct {
	categoryRepo CategoryRepo
}

type CategoryRepo interface {
	GetCategories(context.Context) ([]*domain.Category, error)
}

func NewCategoryService(repo CategoryRepo) *categoryService {
	return &categoryService{categoryRepo: repo}
}

func (c *categoryService) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	return c.categoryRepo.GetCategories(ctx)
}
