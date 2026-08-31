package repository

import (
	"context"
	"database/sql"

	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
	"uuid"
)

// / this would migrated to sqlite package
// DBTX is the query surface shared by *sql.DB and *sql.Tx, so a repo method
// written against it runs standalone or inside a caller's transaction.
type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}


type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// DB is the full handle a repo receives at construction: queries plus the
// ability to open a transaction.
type DB interface {
	DBTX
	Beginner
}

// scanner covers both *sql.Row and *sql.Rows so the scan helpers serve the
// single-row and multi-row queries alike.
type scanner interface {
	Scan(...any) error
}

type ChatRepo struct {
	db DB
}

func NewChatRepo(db DB) *ChatRepo { // repository(for all cases)
	return &ChatRepo{db: db}
}

const createChatQuery = `INSERT into conversations (id) VALUES (?)`

// CreateConversation takes a DBTX so it works both on its own and as one step
// of CreateConversationWithParticipants' transaction.
func (q *ChatRepo) CreateConversation(ctx context.Context, chatID uuid.UUID, db DBTX) error {
	_, err := db.ExecContext(ctx, createChatQuery, chatID.String())
	return sqlite.TranslateError(err)
}

const CreateConversationParticipantsQuery = `INSERT INTO conversation_participants (id, user_id, conversation_id) VALUES (?, ?, ?)`

func (q *ChatRepo) CreateConversationWithParticipants(ctx context.Context, chatID uuid.UUID, idUUID, userIDs []uuid.UUID) error {
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
		_, err = tx.ExecContext(ctx, CreateConversationParticipantsQuery, idUUID[i].String(), userID.String(), chatID.String())
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

func (q *ChatRepo) GetChat(ctx context.Context, usersID []uuid.UUID) (uuid.UUID, error) {
	//                                             b                  a
	row := q.db.QueryRowContext(ctx, getChatQuery, usersID[0].String(), usersID[1].String()) // must be updated later
	var chat uuid.UUID
	err := row.Scan(&chat)
	return chat, err
}

const getMessagesQuery = `
	SELECT id, sender_id, conversation_id, content, created_at FROM messages
	WHERE conversation_id = ?
	ORDER BY created_at DESC, id DESC
	LIMIT ?
	OFFSET ?
`

func (q *ChatRepo) GetMessages(ctx context.Context, chatID uuid.UUID, offset, limit int) ([]domain.MessageOutput, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesQuery, chatID.String(), limit, offset)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	messages := make([]domain.MessageOutput, 0, limit)
	defer rows.Close()
	for rows.Next() {
		var message domain.MessageOutput
		err := rows.Scan(&message.ID, &message.Sender, &message.ChatID, &message.Content, &message.Created_at)
		if err != nil {
			return nil, sqlite.TranslateError(err)
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return messages, nil
}

const sendMessageQuery = `INSERT INTO messages (id, sender_id, content, conversation_id, created_at) VALUES (?, ?, ?, ?, CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER)) RETURNING created_at`

func (q *ChatRepo) SendMessage(ctx context.Context, message domain.MessageOutput) (int64, error) {
	var createdAt int64

	err := q.db.QueryRowContext(ctx, sendMessageQuery, message.ID.String(), message.Sender.String(), message.Content, message.ChatID.String()).Scan(&createdAt)
	if err != nil {
		return -1, sqlite.TranslateError(err)
	}
	return createdAt, nil
}
