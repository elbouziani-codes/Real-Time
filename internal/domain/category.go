package domain

import "realTime/crypto"

type Category struct {
	ID        crypto.UUID `json:"id"`
	Title     string      `json:"title"`
	Icon      string      `json:"icon"`
	CreatedAt int         `json:"created_at"`
}
