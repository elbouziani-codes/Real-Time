package repository

import (
	"context"
	"database/sql"
	"fmt"

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
	_, err := tx.ExecContext(ctx, createChatQuery, chatID.Value)
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
		_, err = tx.ExecContext(ctx, CreateConversationParticipantsQuery, idUUID[i].Value, userID.Value, chatID.Value)
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
/*
CREATE TABLE
	IF NOT EXISTS messages (
		id CHAR(36) PRIMARY KEY,
		sender_id CHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE, -- actuallt this must be reviewed if a user delete whta s the correct practice 
		conversation_id CHAR(36) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours')),
		updated_at INTEGER DEFAULT (unixepoch (CURRENT_TIMESTAMP, '+1 hours'))
	);
	*/
const getMessagesQuery = `
	SELECT id, sender_id, conversation_id, content , created_at FROM messages
	WHERE conversation_id  = ?
	ORDER BY created_at DESC
	LIMIT  ?
	OFFSET ?
`

func (q *ChatRepo) GetMessages(ctx context.Context, chatID crypto.UUID, offset, limit int) ([]domain.MessageOutput, error) {
	fmt.Println(chatID, limit, offset)
	rows, err := q.db.QueryContext(ctx, getMessagesQuery, chatID.Value, limit, offset) // must be updated later
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	messages := make([]domain.MessageOutput, 0, limit)
	defer rows.Close()
	for rows.Next() {
		var message domain.MessageOutput
		err := rows.Scan(&message.ID.Value, &message.Sender.Value, &message.ChatID.Value, &message.Content, &message.Created_at)
		if err != nil { 
			return nil, sqlite.TranslateError(err)
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil{
		return nil, sqlite.TranslateError(err)
	}
	return messages, nil
}

const sendMessageQuery = `INSERT INTO messages (id, sender_id, content, conversation_id) VALUES (?, ?, ?, ?) RETURNING created_at`

func (q *ChatRepo) SendMessage(ctx context.Context, message domain.MessageOutput) (int64, error) {
	var createdAt int64

	err := q.db.QueryRowContext(ctx, sendMessageQuery, message.ID.Value, message.Sender.Value, message.Content, message.ChatID.Value).Scan(&createdAt)
	if err != nil {
		return -1, sqlite.TranslateError(err)
	}
	return createdAt, nil
}
