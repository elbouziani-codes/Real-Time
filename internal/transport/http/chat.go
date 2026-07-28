package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"realTime/crypto"
	"realTime/internal/domain"
)

type HandlerChat struct {
	ChatService ChatService
	UserService UserService
	hub         hub
}

type hub struct {
	lengthAllRead int
}
type ChatService interface {
	IsChatNotFound(error) bool
	CheckRoomChat(context.Context, []crypto.UUID) (crypto.UUID, error)
	GetMessages(context.Context, crypto.UUID, int) ([]domain.MessageOutput, error)
	splitMessages([]domain.MessageOutput, crypto.UUID, crypto.UUID) ([]domain.MessageOutput, []domain.MessageOutput, error)
}
func NewHandleChat(chatService ChatService, userService UserService) *HandlerChat{
	return &HandlerChat{ChatService: chatService , UserService:userService , hub: hub{lengthAllRead:0}}
}
func (chat *HandlerChat) Chat(w http.ResponseWriter , r *http.Request){
	userIDAny := r.Context().Value("user_id")
	userID, _ := userIDAny.(crypto.UUID)
	var chatRoomInput domain.ChatRoomInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chatRoomInput); err != nil{
		http.Error(w,err.Error(), 404)
		return
	}
	err ,chatRoomOrigin :=  chatRoomInput.ValidateAndParse(userID)
	if err != nil{
		http.Error(w,err.Error(), 404)
		return
	}

	err = chat.CheckRoomChat(r.Context(), chatRoomOrigin)
	if err != nil{
		http.Error(w,err.Error(), 404)
		return
	}else if chat.ChatService.IsChatNotFound(err){
		return
	}
	messages, err := chat.ChatService.GetMessages(r.Context(), chatRoomOrigin.ID, chat.hub.lengthAllRead)
	if err != nil{
		http.Error(w,err.Error(), 404)
		return
	}
	messagesMe, messagesFreind, err := chat.ChatService.splitMessages(messages, chatRoomOrigin.Me, chatRoomOrigin.Freind)
	if err != nil{
		http.Error(w,err.Error(), 404)
		return
	}
	
	Response := domain.ChatRoomOutput{ID:chatRoomInput.ID, Me:messagesMe, Freind:messagesFreind}
}

func (chat *HandlerChat) CheckRoomChat(ctx context.Context, chatRoomOrigin *domain.ChatRoomOrigin) error{
	idChat, err := chat.ChatService.CheckRoomChat(ctx, chatRoomOrigin.Particpants)
	if err != nil{
		return err
	}
	if idChat != chatRoomOrigin.ID{
		return domain.Error{Message: "Error in ID Chat NOT Valid", Code: 404}
	}
	return nil
}