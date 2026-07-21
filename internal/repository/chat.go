package repository

import (
	"context"
	"database/sql"

	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
)

// / this would migrated to sqlite package
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type ChatRepo struct {
	db DBTX
}

func NewChatRepo(db DBTX) *ChatRepo { // repository(for all cases)
	return &ChatRepo{db: db}
}

const createChatQuery = `INSERT into conversations (id) VALUES (?)`

func (q *ChatRepo) CreateConversation(ctx context.Context, chatID crypto.UUID, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, createChatQuery, chatID)
	return sqlite.TranslateError(err)
}

const CreateConversationParticipantsQuery = `INSERT INTO conversation_participants (id, user_id, conversation_id) VALUES (?, ?, ?)`

func (q *ChatRepo) CreateConversationWithParticipants(ctx context.Context, chatID crypto.UUID, idUUID, userIDs []crypto.UUID) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	defer tx.Rollback()

	err = q.CreateConversation(ctx, chatID, tx)
	if err != nil {
		return sqlite.TranslateError(err)
	}

	for i, userID := range userIDs {
		_, err = tx.ExecContext(ctx, CreateConversationParticipantsQuery, idUUID[i], userID, chatID)
		if err != nil {
			return sqlite.TranslateError(err)
		}
	}
	return tx.Commit()
}

const getChatQuery = `
	SELECT cp1.conversation_id
	FROM conversation_participants cp1 
	JOIN conversation_participants cp2
	ON cp1.conversation_id = cp2.conversation_id
		WHERE cp1.user_id = ?
		AND cp2.user_id = ?
`



func (q *ChatRepo) GetChat(ctx context.Context, usersID []crypto.UUID) (crypto.UUID, error) {
	//                                             b                  a 
	row := q.db.QueryRowContext(ctx, getChatQuery, usersID[0].Value, usersID[1].Value) // must be updated later
	var chat crypto.UUID
	err := row.Scan(&chat.Value)
	return chat, err
}

const getMessagesQuery = `
	SELECT * FROM messages
	WHERE conversation_id  = ?
	ORDER BY created_at DESC
	OFFSET ?
	LIMIT  ?
`

func (q *ChatRepo) GetMessages(ctx context.Context, chatID string, offset, limit int) ([]domain.Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesQuery, chatID, offset, limit) // must be updated later
	if err != nil {
		return nil, err
	}
	var messages []domain.Message
	defer rows.Close()
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(&message.ID, &message.SenderID, &message.Content)
		if err != nil {
			return nil, sqlite.TranslateError(err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

const sendMessageQuery = `INSERT INTO messages (id, sender_id, content, conversation_id) VALUES (?, ?, ?, ?)`

func (q *ChatRepo) SendMessage(ctx context.Context, message domain.Message) error {
	result, err := q.db.ExecContext(ctx, sendMessageQuery, message.ID, message.SenderID, message.Content, message.ChatID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return sqlite.TranslateError(err)
	}
	if n != 1 { // msut be updated
		return sqlite.TranslateError(err)
	}
	return nil
}
