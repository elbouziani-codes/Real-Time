package domain

import (
	"strings"

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

type MessageInput struct {
	ReceiverID crypto.UUID `json:"receiver_id"`
	Mod        string `json:"mod"`
	ID         crypto.UUID `json:"id"`
	Content    string `json:"content"`
}

func (c MessageInput) ValidMessage() error {
	c.Content = strings.TrimSpace(c.Content)


	if len(c.Content) > 2048 || len(c.Content) == 0 {
		return Error{Message: "length message must be between 1 and 2048  chars", Code: BadFormatCode}
	}

	return nil
}
