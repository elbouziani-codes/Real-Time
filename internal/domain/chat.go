package domain

import (
	"realTime/crypto"
	"slices"
)

type Message struct {
	ID       crypto.UUID
	SenderID crypto.UUID
	ChatID   crypto.UUID
	Content  string
}

type ChatRoomOrigin struct {
	ID          crypto.UUID
	Particpants []crypto.UUID
}

type ChatRoomInput struct {
	ID          string
	Particpants []string
}

func (chatInput *ChatRoomInput) ValidateAndParse(me crypto.UUID) (error, *ChatRoomOrigin){
	if len(chatInput.Particpants) != 2 || (chatInput.Particpants[0] == "" || chatInput.Particpants[1] == "") || chatInput.Particpants[1] == chatInput.Particpants[0]{
		return Error{Message: "error in slice Particpants", Code: 404}, &ChatRoomOrigin{}
	}
	if chatInput.ID == "" {
		return Error{Message: "error in ID chatRoom", Code: 404}, &ChatRoomOrigin{}
	}
	err, chatRoomOrigin := chatInput.ParseUUID()
	if err != nil {
		return err, &ChatRoomOrigin{}
	}
	if slices.Contains(chatRoomOrigin.Particpants, me){
		return Error{Message: "error in ID pirmession for entre this Chatroom", Code: 404}, &ChatRoomOrigin{}
	}
	return nil, chatRoomOrigin
}

func (chatInput *ChatRoomInput) ParseUUID() (error, *ChatRoomOrigin){
	var chatRoomOrigin ChatRoomOrigin
	tab ,err := crypto.ParseUUID(chatInput.ID)
	if err != nil {
		return Error{Message: "error in Parse ID chatRoom", Code: 404}, &ChatRoomOrigin{}
	}
	chatRoomOrigin.ID = tab
	for i := 0; i < len(chatInput.Particpants); i++ {
		tab ,err := crypto.ParseUUID(chatInput.Particpants[i])
		if err != nil {
			return Error{Message: "error in Parse ID chatRoom", Code: 404}, &ChatRoomOrigin{}
		}
		chatRoomOrigin.Particpants = append(chatRoomOrigin.Particpants, tab)
	}
	return nil, &chatRoomOrigin

}