package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"realTime/internal/domain"
	"uuid"
)

type HandlerChat struct {
	ChatService ChatService
	UserService UserService
}

type ChatService interface {
	CheckRoomChat(context.Context, []uuid.UUID) (uuid.UUID, error)
	GetMessages(context.Context, uuid.UUID, int64, uuid.UUID) ([]domain.MessageOutput, error)
}

func NewHandleChat(chatService ChatService, userService UserService) *HandlerChat {
	return &HandlerChat{ChatService: chatService, UserService: userService}
}
func (chat *HandlerChat) Chat(w http.ResponseWriter, r *http.Request) {

	userIDAny := r.Context().Value("user_id")
	userID, _ := userIDAny.(uuid.UUID)

	var chatRoomInput domain.ChatRoomInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chatRoomInput); err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	friendID, err := uuid.Parse(chatRoomInput.Friend)
	if err != nil {
		http.Error(w, "invalid friend id", http.StatusBadRequest)
		return
	}
	if friendID == userID {
		http.Error(w, "cannot open a conversation with yourself", http.StatusBadRequest)
		return
	}

	chatID, err := chat.ChatService.CheckRoomChat(r.Context(), []uuid.UUID{userID, friendID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			json.NewEncoder(w).Encode(&domain.ChatRoomOutput{Messages: []domain.MessageOutput{}})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// A cursor must carry a well-formed id, otherwise it could never match a
	// row; an absent cursor (BeforeAt == 0) is the "newest page" request.
	var beforeID uuid.UUID
	if chatRoomInput.BeforeAt != 0 {
		beforeID, err = uuid.Parse(chatRoomInput.BeforeID)
		if err != nil {
			http.Error(w, "invalid before_id", http.StatusBadRequest)
			return
		}
	}

	messages, err := chat.ChatService.GetMessages(r.Context(), chatID, chatRoomInput.BeforeAt, beforeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Response := domain.ChatRoomOutput{ID: chatID.String(), Messages: messages}
	err = json.NewEncoder(w).Encode(&Response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
