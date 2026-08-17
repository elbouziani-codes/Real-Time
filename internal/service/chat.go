package service

import (
	//	"fmt"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"realTime/crypto"
	"realTime/internal/domain"
	// "realTime/database"
)

type chatService struct {
	UserRepo UserRepo
	ChatRepo ChatRepo
}
type ChatRepo interface {
	CreateConversationWithParticipants(context.Context, crypto.UUID, []crypto.UUID, []crypto.UUID) error
	GetChat(context.Context, []crypto.UUID) (crypto.UUID, error)
	GetMessages(context.Context, crypto.UUID, int, int) ([]domain.MessageOutput, error) 
	SendMessage(context.Context, domain.MessageOutput) (int64, error)
}


func NewChatService(ChatRepo ChatRepo, UserRepo UserRepo) *chatService {
	return &chatService{UserRepo: UserRepo, ChatRepo: ChatRepo}
}

func (c *chatService) IsChatNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (c *chatService) CheckRoomChat(ctx context.Context, userIDs []crypto.UUID) (crypto.UUID, error) {
	chatID, err := c.ChatRepo.GetChat(ctx, userIDs)
	if err != nil {
		return crypto.Nil, err
	}
	return chatID, nil
}

func (c *chatService) CreateRoomChat(ctx context.Context, userIDs []crypto.UUID) (crypto.UUID, error) {
	fmt.Println(userIDs)
	idChat, err := crypto.GenerateUUID()
	if err != nil {
		return crypto.Nil, err
	}
	// func (q *ChatRepo) CreateConversationWithParticipants(ctx context.Context, chatID string, idUUID, userIDs []string) error {

	UUIDs := []crypto.UUID{}
	for range userIDs {

		UUID, err := crypto.GenerateUUID()
		if err != nil {
			return crypto.Nil, err
		}

		UUIDs = append(UUIDs, UUID)
	}

	err = c.ChatRepo.CreateConversationWithParticipants(ctx, idChat, UUIDs, userIDs)
	if err != nil {
		return crypto.Nil, err
	}

	return idChat, nil
}

func (c *chatService) SendMessageRoomChat(ctx context.Context, content string, SenderID, ChatID crypto.UUID) (crypto.UUID, int64, error) {
	idMessage, err := crypto.GenerateUUID()
	if err != nil {
		return crypto.Nil, 0, err
	}

	msg := domain.MessageOutput{ID: idMessage, Content: content, Sender: SenderID, ChatID: ChatID}

	createdAt, err := c.ChatRepo.SendMessage(ctx, msg)
	if err != nil {
		return crypto.Nil, 0, err
	}

	return idMessage, createdAt, nil
}

func (c *chatService) GetMessages(ctx context.Context,chatID crypto.UUID, lengthAllRead int) ([]domain.MessageOutput, error) {
	messages, err := c.ChatRepo.GetMessages(ctx, chatID, lengthAllRead, lengthAllRead+10)
	if err != nil {
		return nil, err
	}
	return messages, err
}
