package domain

import (
	"realTime/crypto"
)

type Message struct {
	ID       crypto.UUID
	SenderID crypto.UUID
	ChatID   crypto.UUID
	Content  string
}

type ChatRoom struct {
	ID          crypto.UUID
	Particpants []crypto.UUID
}
