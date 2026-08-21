package domain

import "realTime/crypto"

type ChatRoomInput struct {
	Friend string `json:"friend"`
	Offset int    `json:"offset"`
}
type MessageOutput struct {
	ID         crypto.UUID `json:"id"`
	Code       int         `json:"code"`
	Sender     crypto.UUID `json:"sender"`
	ChatID     crypto.UUID `json:"chat_id"`
	Content    string      `json:"content"`
	Created_at int64       `json:"created_at"`
}
type ChatRoomOutput struct {
	ID       string          `json:"id"`
	Messages []MessageOutput `json:"me"`
}
