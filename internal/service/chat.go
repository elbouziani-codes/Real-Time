package service

import (
	//	"fmt"
	"context"
	"database/sql"
	"errors"

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
	GetMessages( context.Context, string, int, int) ([]domain.Message, error) 
	SendMessage(context.Context, domain.Message) error
}

func NewChatService(ChatRepo ChatRepo, UserRepo UserRepo) chatService {
	return chatService{UserRepo: UserRepo, ChatRepo: ChatRepo}
}

func (c chatService) IsChatNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (c chatService) CheckRoomChat(ctx context.Context, userIDs []crypto.UUID) (crypto.UUID, error) {
	chatID, err := c.ChatRepo.GetChat(ctx, userIDs)
	if err != nil {
		return crypto.Nil, err
	}
	return chatID, nil
}

func (c chatService) CreateRoomChat(ctx context.Context, userIDs []crypto.UUID) (crypto.UUID, error) {
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

func (c chatService) SendMessageRoomChat(ctx context.Context, content string, SenderID, ChatID crypto.UUID) (crypto.UUID, error) {
	idMessage, err := crypto.GenerateUUID()
	if err != nil {
		return crypto.Nil, err
	}

	msg := domain.Message{ID: idMessage, Content: content, SenderID: SenderID, ChatID: ChatID}

	err = c.ChatRepo.SendMessage(ctx, msg)
	if err != nil {
		return crypto.Nil, err
	}

	return idMessage, nil
}
