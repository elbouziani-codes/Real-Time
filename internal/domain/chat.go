package domain

import "uuid"

type ChatRoomInput struct {
	Friend string `json:"friend"`
	Offset int    `json:"offset"`
}
type MessageOutput struct {
	ID         uuid.UUID `json:"id"`
	Code       int       `json:"code"`
	Sender     uuid.UUID `json:"sender"`
	ChatID     uuid.UUID `json:"chat_id"`
	Content    string    `json:"content"`
	Created_at int64     `json:"created_at"`
}

type ChatRoomOutput struct {
	ID       string          `json:"id"`
	Messages []MessageOutput `json:"me"`
}
