package repository

import (
	"context"
	"database/sql"

	"realTime/domain"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type ChatRepo struct {
	db DBTX
}

func NewChatRepo(db DBTX) *ChatRepo { // repository(for all cases)
	return &ChatRepo{db: db}
}

const createChatQuery = `
	INSERT into conversations (id) VALUES (?) RETURNING * ;
`

func (q *ChatRepo) CreateChat(ctx context.Context, chatID []byte) (domain.ChatRoom, error) {
	row := q.db.QueryRowContext(ctx, createChatQuery, chatID)
	var chat domain.ChatRoom
	err := row.Scan(&chat.ID)
	return chat, err
}

const getChatQuery = `
	SELECT cp1.conversation_id 		
	FROM conversation_participants cp1
	JOIN conversation_participants cp2
	ON cp1.conversation_id = cp2.conversation_id
		WHERE cp1.user_id = ? 
		AND cp2.user_id = ?  
`

func (q *ChatRepo) GetChat(ctx context.Context, users []domain.User) (domain.ChatRoom, error) {
	row := q.db.QueryRowContext(ctx, getChatQuery, users[0].ID, users[1].ID) // must be updated later
	var chat domain.ChatRoom
	err := row.Scan(&chat.ID)
	return chat, err
}

const getMessagesQuery = `
	SELECT * FROM messages 
	WHERE conversation_id  = ? 
	ORDER BY created_at DESC
	OFFSET ? 
	LIMIT  ? 	
`

func (q *ChatRepo) GetMessages(ctx context.Context, chatID []byte) ([]domain.Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesQuery, chatID) // must be updated later
	if err != nil {
		return nil, err
	}
	var messages []domain.Message
	defer rows.Close()
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(&message.ID, &message.Sender, &message.Content)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

const sendMessageQuery = `INSERT INTO messages (id, sender_id, value) VALUES (?, ?, ?)`

func (q *ChatRepo) SendMessage(ctx context.Context, message domain.Message) error {
	row, err := q.db.ExecContext(ctx, sendMessageQuery, message.Content)
	if err != nil {
		return err
	}
	n, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return err
	}
	return nil
}
