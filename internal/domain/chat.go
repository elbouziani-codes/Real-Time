package domain

import (
	"context"
)

type Message struct {
	ID string
	Sender string
	Content string	
}

type ChatRoom struct {
	ID string
	Particpants []User	
}

type RepoChat interface {
	GetChat(context.Context, [][]User) (ChatRoom, error)
	CreateChat(context.Context, string) (ChatRoom, error)
	GetMessage(context.Context) ([]Message, error)
	AddUsersToChat(context.Context, []User) (error)
} 
