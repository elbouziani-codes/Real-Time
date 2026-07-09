package domain


type Message struct {
	ID []byte
	Sender []byte
	Content string	
}

type ChatRoom struct {
	ID []byte
	Particpants []User	
}

type RepoChat interface {
	GetChat(context.Context, [][]User) (ChatRoom, error)
	CreateChat(context.Context, []byte) (ChatRoom, error)
	GetMessage(context.Context) ([]Message, error)
	AddUsersToChat(context.Context, []User) (error)
} 
