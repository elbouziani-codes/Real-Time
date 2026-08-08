package repository

import (
	"context"
	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
)

type categoryRepo struct {
	db DBTX
}

func NewCategoryRepo(db DBTX) *categoryRepo {
	return &categoryRepo{db: db}
}

func scanCategory(row scanner) (*domain.Category, error) {
	category := domain.Category{}
	err := row.Scan(
		&category.ID.Value,
		&category.Title,
		&category.Icon,
		&category.CreatedAt)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return &category, nil
}

const getCategoriesQuery = `
SELECT C.id, C.title, C.icon, C.created_at
FROM categories C
ORDER BY C.title ASC
`

func (c *categoryRepo) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	rows, err := c.db.QueryContext(ctx, getCategoriesQuery)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	defer rows.Close()
	var categories []*domain.Category
	for rows.Next() {
		category, err := scanCategory(rows)
		if err != nil {
			return nil, sqlite.TranslateError(err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return categories, nil
}

const getCategoryQuery = `
SELECT C.id, C.title, C.icon, C.created_at
FROM categories C
WHERE C.id = ?
`

func (c *categoryRepo) GetCategory(ctx context.Context, categoryID crypto.UUID) (*domain.Category, error) {
	row := c.db.QueryRowContext(ctx, getCategoryQuery, categoryID.Value)
	return scanCategory(row)
}
