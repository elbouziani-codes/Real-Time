package handler

import (
	"context"
	"net/http"
	"sort"
	"sync"

	"realTime/internal/domain"

	"github.com/gorilla/websocket"
)

type handlerWs struct {
	svcChat ChatService
	svcUser UserService
	hub     *hub
}

type hub struct {
	mu      sync.RWMutex
	clients map[string][]*Client
}
type UserService interface {
	GetByID(context.Context, string) (domain.User, error)
}

type ChatService interface {
	IsChatNotFound(error) bool
	CheckRoomChat(context.Context, []string) (string, error)
	CreateRoomChat(context.Context, []string) (string, error)
	SendMessageRoomChat(context.Context, string, string, string) (string, error)
	ValidMessage(context.Context, string, string) error
}

func NewHandleWs(svcChat ChatService, svcUser UserService) *handlerWs {
	return &handlerWs{
		svcChat: svcChat,
		svcUser: svcUser,
		hub: &hub{
			clients: make(map[string][]*Client),
		},
	}
}

type MessageInput struct {
    ReceiverID string `json:"receiver_id"`
    Content string `json:"content"`
}

type MessageOutput struct {
	ID      string `json:"id"`
	Code    int    `json:"code"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type Client struct {
	id     int
	UserId string
	conn   *websocket.Conn
	mu      sync.Mutex
}

var allConn = make(map[string][]*Client)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wss handlerWs) chatWs(w http.ResponseWriter, r *http.Request) {
	client := wss.AddClient(w,r)
	
	go wss.engineMessages(r.Context(), client)
}

func (wss handlerWs) engineMessages(ctx context.Context, client *Client) {
	defer wss.DeleteClient(client)
	for {
		var MessageInput MessageInput
		err := client.conn.ReadJSON(&MessageInput)
		if err != nil {
			return
		}
		err = wss.svcChat.ValidMessage(ctx, MessageInput.ReceiverID, MessageInput.Content)
		if err != nil {
			responceWrite("", 404, err.Error(), client.UserId, client)
			continue
		}

		userIDs := []string{client.UserId, MessageInput.ReceiverID}
		sort.Strings(userIDs)
		chatId, err := wss.GetRoom(ctx, userIDs)
		if chatId == "-99" && err != nil {
			responceWrite("", 404, err.Error(), client.UserId, client)
			continue
		}
		msgId, err := wss.svcChat.SendMessageRoomChat(ctx, MessageInput.Content, client.UserId, chatId)
		wss.hub.mu.RLock()
			all_receiver := make([]*Client, len(wss.hub.clients[MessageInput.ReceiverID]))
			copy(all_receiver, wss.hub.clients[MessageInput.ReceiverID])

		wss.hub.mu.RUnlock()
		if all_receiver != nil {
			for _, receiver := range all_receiver {
				responceWrite(msgId, 200, MessageInput.Content, MessageInput.ReceiverID, receiver)
			}
		}
	}
}

func responceWrite(id string, codeError int, content string, SenderID string, receiver *Client) {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	MessageOutput := MessageOutput{ID: id, Code: codeError, Content: content, Sender: SenderID}
	receiver.conn.WriteJSON(&MessageOutput)
}

func (wss handlerWs) GetRoom(ctx context.Context, userIDs []string) (string, error) {
	idChat, err := wss.svcChat.CheckRoomChat(ctx, userIDs)
	if wss.svcChat.IsChatNotFound(err) {
		idChat, err = wss.svcChat.CreateRoomChat(ctx, userIDs)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return idChat, err
}

func (wss handlerWs) DeleteClient(client *Client) {
	wss.hub.mu.Lock()
	defer wss.hub.mu.Unlock()
	userConnections, exists := wss.hub.clients[client.UserId]
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
		delete(wss.hub.clients, client.UserId)
		return
	}

	wss.hub.clients[client.UserId] = userConnections
}


func (wss handlerWs) AddClient(w http.ResponseWriter , r *http.Request) (*Client){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return &Client{}
	}


	userAny := r.Context().Value("me")
	user := userAny.(domain.User)
	client := &Client{
		UserId: user.ID,
		conn:   conn,
	}


	wss.hub.mu.Lock()

	if userConnections := wss.hub.clients[client.UserId]; userConnections == nil {
		client.id = 1
		var clients []*Client
		wss.hub.clients[client.UserId] = clients
	} else {
		num := len(wss.hub.clients[client.UserId])
		client.id = num + 1
	}
	wss.hub.clients[client.UserId] = append(wss.hub.clients[client.UserId], client)
	wss.hub.mu.Unlock()
	return client
}
