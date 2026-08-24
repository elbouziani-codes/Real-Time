package domain

import "uuid"

// ChatRoomInput pages a conversation history. BeforeAt/BeforeID name the
// oldest message the client already holds (keyset paging); both zero mean
// "the newest page", so an opening request carries no cursor.
type ChatRoomInput struct {
	Friend   string `json:"friend"`
	BeforeAt int64  `json:"before_at"`
	BeforeID string `json:"before_id"`
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
