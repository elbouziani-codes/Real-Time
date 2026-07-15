package domain

type Message struct {
	ID      string
	Sender  string
	Content string
}

type ChatRoom struct {
	ID          string
	Particpants []User
}
