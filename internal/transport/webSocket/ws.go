package ws

import (
	"context"
	"net/http"
	"sync"

	"realTime/crypto"
	"realTime/internal/domain"

	"github.com/gorilla/websocket"
)

type HandlerWs struct {
	svcChat ChatService
	svcUser UserService
	hub     *hub
}

type hub struct {
	mu      sync.RWMutex
	clients map[crypto.UUID][]*Client
}
type UserService interface {
	GetByID(context.Context, crypto.UUID) (domain.User, error)
}

type ChatService interface {
	IsChatNotFound(error) bool
	CheckRoomChat(context.Context, []crypto.UUID) (crypto.UUID, error)
	CreateRoomChat(context.Context, []crypto.UUID) (crypto.UUID, error)
	SendMessageRoomChat(context.Context, string, crypto.UUID, crypto.UUID) (crypto.UUID, error)
}

func NewHandleWs(svcChat ChatService, svcUser UserService) HandlerWs {
	return HandlerWs{
		svcChat: svcChat,
		svcUser: svcUser,
		hub: &hub{
			clients: make(map[crypto.UUID][]*Client),
		},
	}
}

type MessageOutput struct {
	ID      crypto.UUID `json:"id"`
	Code    int         `json:"code"`
	Sender  crypto.UUID `json:"sender"`
	Content string      `json:"content"`
}

type Client struct {
	id     int
	userId crypto.UUID
	conn   *websocket.Conn
	mu     sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wss HandlerWs) ChatWs(w http.ResponseWriter, r *http.Request) {
	client := wss.AddClient(w, r)

	go wss.engineMessages(r.Context(), client)
}

func (wss HandlerWs) engineMessages(ctx context.Context, client *Client) {
	defer wss.DeleteClient(client)
	for {
		var messageInput domain.MessageInput
		err := client.conn.ReadJSON(&messageInput)
		if err != nil {
			return
		}
		err = messageInput.ValidMessage()
		// wss.CheckReceiver()
		if err != nil {
			responceWrite(crypto.Nil, 404, err.Error(), client.userId, client)
			continue
		}

		_, err = wss.svcUser.GetByID(ctx, messageInput.ReceiverID)
		if err != nil {
			responceWrite(crypto.Nil, 404, err.Error(), client.userId, client)
			continue
		}

		userIDs := []crypto.UUID{client.userId, messageInput.ReceiverID}
		chatId, err := wss.GetRoom(ctx, userIDs)
		if err != nil {
			responceWrite(crypto.Nil, 404, err.Error(), client.userId, client)
			continue
		}
		msgId, err := wss.svcChat.SendMessageRoomChat(ctx, messageInput.Content, client.userId, chatId)
		wss.hub.mu.RLock()
		all_receiver := make([]*Client, len(wss.hub.clients[messageInput.ReceiverID]))
		copy(all_receiver, wss.hub.clients[messageInput.ReceiverID])

		wss.hub.mu.RUnlock()
		if all_receiver != nil {
			for _, receiver := range all_receiver {
				responceWrite(msgId, 200, messageInput.Content, messageInput.ReceiverID, receiver)
			}
		}
	}
}

func responceWrite(id crypto.UUID, codeError int, content string, SenderID crypto.UUID, receiver *Client) {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	MessageOutput := MessageOutput{ID: id, Code: codeError, Content: content, Sender: SenderID}
	receiver.conn.WriteJSON(&MessageOutput)
}

func (wss HandlerWs) GetRoom(ctx context.Context, userIDs []crypto.UUID) (crypto.UUID, error) {
	idChat, err := wss.svcChat.CheckRoomChat(ctx, userIDs)
	if wss.svcChat.IsChatNotFound(err) {
		idChat, err = wss.svcChat.CreateRoomChat(ctx, userIDs)
		if err != nil {
			return crypto.Nil, err
		}
	} else if err != nil {
		return crypto.Nil, err
	}
	return idChat, err
}

func (wss HandlerWs) DeleteClient(client *Client) {
	wss.hub.mu.Lock()
	defer wss.hub.mu.Unlock()
	userConnections, exists := wss.hub.clients[client.userId]
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
		delete(wss.hub.clients, client.userId)
		return
	}

	wss.hub.clients[client.userId] = userConnections
}

func (wss HandlerWs) AddClient(w http.ResponseWriter, r *http.Request) *Client {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return &Client{}
	}

	userIDAny := r.Context().Value("user_id")

	userID := userIDAny.(crypto.UUID)
	client := &Client{
		userId: userID,
		conn:   conn,
	}

	wss.hub.mu.Lock()

	if userConnections := wss.hub.clients[client.userId]; userConnections == nil {
		client.id = 1
		var clients []*Client
		wss.hub.clients[client.userId] = clients
	} else {
		num := len(wss.hub.clients[client.userId])
		client.id = num + 1
	}
	wss.hub.clients[client.userId] = append(wss.hub.clients[client.userId], client)
	wss.hub.mu.Unlock()
	return client
}
