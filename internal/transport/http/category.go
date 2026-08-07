package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"realTime/crypto"
	"realTime/internal/domain"
)

type categoryService interface {
	GetCategories(context.Context) ([]*domain.Category, error)
	GetCategory(context.Context, crypto.UUID) (*domain.Category, error)
}

type CategoryHandler struct {
	categorySvc categoryService
}

func NewCategoryHandler(svc categoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc: svc}
}

func (c *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.categorySvc.GetCategories(r.Context())
	if err != nil {
		Error(err, w)
		return
	}

	if categories == nil {
		categories = []*domain.Category{}
	}

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "uknown error", http.StatusInternalServerError)
		return
	}
}
