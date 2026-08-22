package domain

import "uuid"

type Category struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Icon      string    `json:"icon"`
	CreatedAt int       `json:"created_at"`
}
