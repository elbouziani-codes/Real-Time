package domain

import (
	"realTime/crypto"
)

type ChatRoomOrigin struct {
	ID          crypto.UUID
	Me 			crypto.UUID	
	Freind 		crypto.UUID
}

type ChatRoomInput struct {
	ID          string    	`json:"id"`
	Me 			string		`json:"me"`
	Freind 		string 		`json:"freind"`
	Offset		int			`json:"offset"`
}
type MessageOutput struct {
	ID      	crypto.UUID 	`json:"id"`
	Code    	int         	`json:"code"`
	Sender  	crypto.UUID 	`json:"sender"`
	ChatID   	crypto.UUID 	`json:"chat_id"`
	Content 	string       	`json:"content"`
	Created_at 	int64			`json:"created_at"`
}
type ChatRoomOutput struct {
	ID          string    		`json:"id"`
	Messages 	[]MessageOutput `json:"me"`

}
func (chatInput *ChatRoomInput) ValidateAndParse(me crypto.UUID) (error, *ChatRoomOrigin){
	if (chatInput.Me == ""  || chatInput.Freind == "") || chatInput.Me == chatInput.Freind{
		return Error{Message: "error in slice Particpants", Code: 404}, &ChatRoomOrigin{}
	}
	if chatInput.ID == "" {
		return Error{Message: "error in ID chatRoom", Code: 404}, &ChatRoomOrigin{}
	}
	err, chatRoomOrigin := chatInput.ParseUUID()
	if err != nil {
		return err, &ChatRoomOrigin{}
	}
	if chatRoomOrigin.Me !=  me{
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
	
	tab ,err = crypto.ParseUUID(chatInput.Me)
	if err != nil {
		return Error{Message: "error in Parse ID chatRoom", Code: 404}, &ChatRoomOrigin{}
	}
	chatRoomOrigin.Me = tab

	tab ,err = crypto.ParseUUID(chatInput.Freind)
	if err != nil {
		return Error{Message: "error in Parse ID chatRoom", Code: 404}, &ChatRoomOrigin{}
	}
	chatRoomOrigin.Freind = tab


	return nil, &chatRoomOrigin

}