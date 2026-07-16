package domain

type Message struct {
    ID      string
    SenderID  string
    ChatID string
    Content string
}
type ChatRoom struct {
	ID          string
	Particpants []string
}
