package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"realTime/crypto"
	"realTime/internal/domain"
)

type HandlerChat struct {
	ChatService ChatService
	UserService UserService
}

type ChatService interface {
	CheckRoomChat(context.Context, []crypto.UUID) (crypto.UUID, error)
	GetMessages(context.Context, crypto.UUID, int) ([]domain.MessageOutput, error)
}

func NewHandleChat(chatService ChatService, userService UserService) *HandlerChat {
	return &HandlerChat{ChatService: chatService, UserService: userService}
}
func (chat *HandlerChat) Chat(w http.ResponseWriter, r *http.Request) {

	userIDAny := r.Context().Value("user_id")
	userID, _ := userIDAny.(crypto.UUID)

	var chatRoomInput domain.ChatRoomInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chatRoomInput); err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	friendID, err := crypto.ParseUUID(chatRoomInput.Friend)
	if err != nil {
		http.Error(w, "invalid friend id", http.StatusBadRequest)
		return
	}
	if friendID == userID {
		http.Error(w, "cannot open a conversation with yourself", http.StatusBadRequest)
		return
	}

	chatID, err := chat.ChatService.CheckRoomChat(r.Context(), []crypto.UUID{userID, friendID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			json.NewEncoder(w).Encode(&domain.ChatRoomOutput{Messages: []domain.MessageOutput{}})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	messages, err := chat.ChatService.GetMessages(r.Context(), chatID, chatRoomInput.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Response := domain.ChatRoomOutput{ID: chatID.Value.String(), Messages: messages}
	err = json.NewEncoder(w).Encode(&Response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
