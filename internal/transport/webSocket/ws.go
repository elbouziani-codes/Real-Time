package handler

import (
	"context"
	"net/http"

	"realTime/internal/domain"

	"github.com/gorilla/websocket"
)

type handlerWs struct {
	svcChat ChatService
	svcUser UserService
}

type UserService interface {
	GetByID(context.Context, string) (domain.User, error)
}

type ChatService interface {
	CheckRoomChat(context.Context, []string) (string, error)
	CreateRoomChat(context.Context, []string) (string, error)
	SendMessageRoomChat(context.Context, string, string, string) (string, error)
	ValidMessage(context.Context, string, string) error
}

func NewHandleWs(svcChat ChatService, svcUser UserService) *handlerWs {
	return &handlerWs{
		svcChat: svcChat,
		svcUser: svcUser,
	}
}

type MessageInput struct {
	receiver_id string
	content     string
}

type MessageOutput struct {
	id      string
	code    int
	Sender  string
	content string
}

type Client struct {
	id     int
	UserId string
	conn   *websocket.Conn
}

var allConn = make(map[string][]*Client)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wss handlerWs) chatWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	userAny := r.Context().Value("me")
	user := userAny.(domain.User)
	client := &Client{
		UserId: user.ID,
		conn:   conn,
	}

	if userConnections := allConn[client.UserId]; userConnections == nil {
		client.id = 1
		var clients []*Client
		allConn[client.UserId] = clients
	} else {
		num := len(allConn[client.UserId])
		client.id = num + 1
	}
	allConn[client.UserId] = append(allConn[client.UserId], client)
	go wss.engineMessages(r.Context(), client)
}

func (wss handlerWs) engineMessages(ctx context.Context, client *Client) {
	for {
		defer DeleteClient(client)
		var MessageInput MessageInput
		err := client.conn.ReadJSON(&MessageInput)
		if err != nil {
			MessageOutput := MessageOutput{id: "", code: 404, content: "Error in Read Missage", Sender: client.UserId}
			client.conn.WriteJSON(&MessageOutput)
			continue
		}
		err = wss.svcChat.ValidMessage(ctx, MessageInput.receiver_id, MessageInput.content)
		if err != nil {
			MessageOutput := MessageOutput{id: "", code: 404, content: err.Error(), Sender: client.UserId}
			client.conn.WriteJSON(&MessageOutput)
			continue
		}
		userIDs := []string{client.UserId, MessageInput.receiver_id}
		chatId := wss.GetRoom(ctx, userIDs)
		msgId, err := wss.svcChat.SendMessageRoomChat(ctx, MessageInput.content, client.UserId, chatId)

		if all_receiver := allConn[MessageInput.receiver_id]; all_receiver != nil {
			for _, receiver := range all_receiver {
				MessageOutput := MessageOutput{id: msgId, code: 200, content: MessageInput.content, Sender: MessageInput.receiver_id}
				receiver.conn.WriteJSON(&MessageOutput)
			}
		}

	}
}

func (wss handlerWs) GetRoom(ctx context.Context, userIDs []string) string {
	idChat, err := wss.svcChat.CheckRoomChat(ctx, userIDs)
	if err.Error() == "error not find" && idChat == "-99" {
		idChat, err = wss.svcChat.CreateRoomChat(ctx, userIDs)
	}

	return idChat
}

func DeleteClient(client *Client) {
    userConnections, exists := allConn[client.UserId]

    if !exists {
        return
    }

    for i, c := range userConnections {
        if c.id == client.id {
            userConnections = append(userConnections[:i], userConnections[i+1:]...)
            break
        }
    }

    if len(userConnections) == 0 {
        delete(allConn, client.UserId)
        return
    }

    allConn[client.UserId] = userConnections
}
