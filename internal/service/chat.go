package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"realTime/internal/domain"
	"uuid"
)

type chatService struct {
	UserRepo UserRepo
	ChatRepo ChatRepo
}
type ChatRepo interface {
	CreateConversationWithParticipants(context.Context, uuid.UUID, []uuid.UUID, []uuid.UUID) error
	GetChat(context.Context, []uuid.UUID) (uuid.UUID, error)
	GetMessages(context.Context, uuid.UUID, int64, uuid.UUID, int) ([]domain.MessageOutput, error)
	SendMessage(context.Context, domain.MessageOutput) (int64, error)
}

func NewChatService(ChatRepo ChatRepo, UserRepo UserRepo) *chatService {
	return &chatService{UserRepo: UserRepo, ChatRepo: ChatRepo}
}

func (c *chatService) IsChatNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (c *chatService) CheckRoomChat(ctx context.Context, userIDs []uuid.UUID) (uuid.UUID, error) {
	chatID, err := c.ChatRepo.GetChat(ctx, userIDs)
	if err != nil {
		return uuid.Nil(), err
	}
	return chatID, nil
}

func (c *chatService) CreateRoomChat(ctx context.Context, userIDs []uuid.UUID) (uuid.UUID, error) {
	fmt.Println(userIDs)
	idChat := uuid.NewV4()

	UUIDs := []uuid.UUID{}
	for range userIDs {
		UUIDs = append(UUIDs, uuid.NewV4())
	}

	err := c.ChatRepo.CreateConversationWithParticipants(ctx, idChat, UUIDs, userIDs)
	if err != nil {
		return uuid.Nil(), err
	}

	return idChat, nil
}

func (c *chatService) SendMessageRoomChat(ctx context.Context, content string, SenderID, ChatID uuid.UUID) (uuid.UUID, int64, error) {
	idMessage := uuid.NewV4()

	msg := domain.MessageOutput{ID: idMessage, Content: content, Sender: SenderID, ChatID: ChatID}

	createdAt, err := c.ChatRepo.SendMessage(ctx, msg)
	if err != nil {
		return uuid.Nil(), 0, err
	}

	return idMessage, createdAt, nil
}

<<<<<<< HEAD
// pageSize is how many messages one history request returns.
const pageSize = 10

func (c *chatService) GetMessages(ctx context.Context, chatID uuid.UUID, beforeAt int64, beforeID uuid.UUID) ([]domain.MessageOutput, error) {
	messages, err := c.ChatRepo.GetMessages(ctx, chatID, beforeAt, beforeID, pageSize)
=======
func (c *chatService) GetMessages(ctx context.Context, chatID uuid.UUID, offset int) ([]domain.MessageOutput, error) {
	messages, err := c.ChatRepo.GetMessages(ctx, chatID, offset, 10)
>>>>>>> melbouzi
	if err != nil {
		return nil, err
	}
	return messages, err
}
