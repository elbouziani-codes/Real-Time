package service

import (
	//	"fmt"
	"context"
	"database/sql"
	"errors"
	"strings"

	"realTime/crypto"
	"realTime/internal/domain"
	// "realTime/database"
)

type chatService struct {
	UserRepo UserRepo
	ChatRepo ChatRepo
}
type ChatRepo interface {
	CreateChat(context.Context, string) error
	GetChat(context.Context, []string) (string, error)
	GetMessages(context.Context, string) ([]domain.Message, error)
	SendMessage( context.Context, domain.Message) error
}

func NewService(ChatRepo ChatRepo, UserRepo UserRepo) chatService {
	return chatService{UserRepo: UserRepo, ChatRepo: ChatRepo}
}

func (c chatService) ValidMessage(ctx context.Context, receiver_id, Content string) error {
	receiver_id = strings.TrimSpace(receiver_id)
	Content = strings.TrimSpace(Content)
	if receiver_id == "" {
		return errors.New("error empty receiver_id")
	}
	_, err := c.UserRepo.GetByID(ctx, receiver_id)
	if err != nil {
		return err
	}
	if len(Content) > 2048 {
		return errors.New("error length message 2048>")
	}
	return nil
}



func (c chatService) CheckRoomChat(ctx context.Context, userIDs []string) (string, error) {
	chatID, err := c.ChatRepo.GetChat(ctx, userIDs)
	if err == sql.ErrNoRows {
		return "-99", errors.New("error not find")
	}
	return chatID, nil
}

func (c chatService) CreateRoomChat(ctx context.Context, userIDs []string) (string, error) {
	
	idChat, err := crypto.GenerateUUID()
	err = c.ChatRepo.CreateChat(ctx, idChat)
	if err != nil {
		return "-99", err
	}
	
	return idChat, err
}

func (c chatService) SendMessageRoomChat(ctx context.Context ,content,SenderID,ChatID string) ( string ,error) {
	idMessage, err := crypto.GenerateUUID()
	if err != nil {
		return "", err
	}

	msg := domain.Message{ID:idMessage, Content: content , SenderID:SenderID , ChatID: ChatID}

	err = c.ChatRepo.SendMessage(ctx, msg)
	if err != nil {
		return "", err
	}

	return idMessage ,nil
}

