package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"uuid"

	"realTime/internal/domain"

	"github.com/gorilla/websocket"
)

/**
 * 1 Close WebSocket Connection and Disable Chat Page
 * 2 typing
 * 3 add client
 * 4 closed webSocket client
 */
type HandlerWs struct {
	svcChat ChatService
	svcUser UserService
	hub     *hub
}

type hub struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]*Client
}
type UserService interface {
	GetByID(context.Context, uuid.UUID) (domain.User, error)
}

type ChatService interface {
	IsChatNotFound(error) bool
	CheckRoomChat(context.Context, []uuid.UUID) (uuid.UUID, error)
	CreateRoomChat(context.Context, []uuid.UUID) (uuid.UUID, error)
	SendMessageRoomChat(context.Context, string, uuid.UUID, uuid.UUID) (uuid.UUID, int64, error)
}

func NewHandleWs(svcChat ChatService, svcUser UserService) *HandlerWs {
	return &HandlerWs{
		svcChat: svcChat,
		svcUser: svcUser,
		hub: &hub{
			clients: make(map[uuid.UUID]*Client),
		},
	}
}

type Client struct {
	userId uuid.UUID
	conn   *websocket.Conn
	mu     sync.Mutex
}

func (c *Client) WriteJSON(v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(v)
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.Close()
}

var upgrader = websocket.Upgrader{}

func (wss *HandlerWs) ChatWs(w http.ResponseWriter, r *http.Request) {
	client, err := wss.AddClient(w, r)
	if err != nil {
		return
	}

	go wss.engineMessages(context.Background(), client)
}

func (wss *HandlerWs) engineMessages(ctx context.Context, sender *Client) {
	defer func() {
		wss.DeleteClient(sender)
	}()
	for {
		err, wsRequest := wss.readInputMessage(ctx, sender)
		if err != nil {
			if err.Error() == "Close" {
				return
			}
			continue
		}
		chatId, err := wss.GetRoom(ctx, sender.userId, wsRequest.Destination)
		if err != nil {
			wss.responseWrite(uuid.Nil(), 404, err.Error(), sender.userId, sender, chatId, 0)
			continue
		}

		if wsRequest.RequestType == "typing" {
			wss.sendTypingReceiver(wsRequest.Destination, sender.userId, chatId)
			continue
		}

		messageID, createAt, err := wss.svcChat.SendMessageRoomChat(ctx, wsRequest.Content, sender.userId, chatId)
		if err != nil {
			wss.responseWrite(uuid.Nil(), 404, err.Error(), sender.userId, sender, chatId, createAt)
			continue
		}

		wss.sendReceiver(wsRequest, messageID, sender.userId, chatId, createAt)
	}
}

func (wss *HandlerWs) sendReceiver(wsRequest *domain.WsParsedRequest, messageID uuid.UUID, senderID uuid.UUID, chatId uuid.UUID, create_at int64) {
	wss.hub.mu.RLock()
	receiver := wss.hub.clients[wsRequest.Destination]
	wss.hub.mu.RUnlock()
	if receiver != nil {
		wss.responseWrite(messageID, 200, wsRequest.Content, senderID, receiver, chatId, create_at)
	}
}

func (wss *HandlerWs) sendTypingReceiver(receiverID uuid.UUID, senderID uuid.UUID, chatId uuid.UUID) {
	wss.hub.mu.RLock()
	receiver := wss.hub.clients[receiverID]
	wss.hub.mu.RUnlock()
	if receiver != nil {
		wss.responseWrite(uuid.Nil(), 2, "typing", senderID, receiver, chatId, 0)
	}
}

func (wss *HandlerWs) readInputMessage(ctx context.Context, client *Client) (error, *domain.WsParsedRequest) {
	var wsRequest domain.WsRequest
	err := client.conn.ReadJSON(&wsRequest)
	if err != nil {
		return errors.New("Close"), &domain.WsParsedRequest{}
	}

	wsRequestParse, err := wsRequest.ValidRequest()
	if err != nil {
		wss.responseWrite(uuid.Nil(), 404, err.Error(), client.userId, client, uuid.Nil(), 0)

		return err, &domain.WsParsedRequest{}
	}

	_, err = wss.svcUser.GetByID(ctx, wsRequestParse.Destination)
	if err != nil {
		wss.responseWrite(uuid.Nil(), 404, err.Error(), client.userId, client, uuid.Nil(), 0)
		return err, &domain.WsParsedRequest{}
	}

	if client.userId.String() == wsRequestParse.Destination.String() {
		err := domain.Error{Message: "Error in sender == Destination", Code: domain.ConflictCode}
		wss.responseWrite(uuid.Nil(), 404, err.Error(), client.userId, client, uuid.Nil(), 0)
		return err, &domain.WsParsedRequest{}
	}

	return nil, wsRequestParse
}

func (wss *HandlerWs) responseWrite(id uuid.UUID, codeError int, content string, senderID uuid.UUID, receiver *Client, chatID uuid.UUID, createdAt int64) {
	messageOutput := domain.MessageOutput{ID: id, Code: codeError, Content: content, Sender: senderID, ChatID: chatID, Created_at: createdAt}
	if err := receiver.WriteJSON(&messageOutput); err != nil {
		wss.DeleteClient(receiver)
	}
}

func (wss *HandlerWs) GetRoom(ctx context.Context, userId uuid.UUID, ReceiverID uuid.UUID) (uuid.UUID, error) {
	userIDs := []uuid.UUID{userId, ReceiverID}
	idChat, err := wss.svcChat.CheckRoomChat(ctx, userIDs)
	if wss.svcChat.IsChatNotFound(err) {
		idChat, err = wss.svcChat.CreateRoomChat(ctx, userIDs)
		if err != nil {
			return uuid.Nil(), err
		}
	} else if err != nil {
		return uuid.Nil(), err
	}
	return idChat, err
}

func (wss *HandlerWs) DeleteClient(client *Client) {
	wss.hub.mu.Lock()

	current, exists := wss.hub.clients[client.userId]

	if !exists || current != client {
		wss.hub.mu.Unlock()
		return
	}

	delete(wss.hub.clients, client.userId)

	wss.hub.mu.Unlock()
	client.Close()
	wss.seedAllClient(client, false)
}

func (wss *HandlerWs) AddClient(w http.ResponseWriter, r *http.Request) (*Client, error) {
	userIDAny := r.Context().Value("user_id")
	userID := userIDAny.(uuid.UUID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	newClient := &Client{
		userId: userID,
		conn:   conn,
	}

	var oldClient *Client
	wss.hub.mu.Lock()
	oldClient = wss.hub.clients[userID]
	wss.hub.clients[userID] = newClient
	wss.hub.mu.Unlock()

	if oldClient != nil {
		
		oldClient.WriteJSON(&domain.MessageOutput{ID: uuid.Nil(), Code: 1, Content: "Close WebSocket Connection and Disable Chat Page", Sender: oldClient.userId, ChatID: uuid.Nil(), Created_at: 0})
		oldClient.Close()
		wss.seedAllClient(oldClient, false)
	}

	wss.seedAllClient(newClient, true)

	return newClient, nil
}

func (wss *HandlerWs) seedAllClient(Me *Client, addNewClient bool) {
	wss.hub.mu.RLock()
	allListClient := make(map[uuid.UUID]*Client, len(wss.hub.clients))

	for uuid, client := range wss.hub.clients {
		allListClient[uuid] = client
	}

	wss.hub.mu.RUnlock()

	allClient := []string{}

	for id, client := range allListClient {
		if Me.userId == id {
			continue
		}
	
		if addNewClient {
			allClient = append(allClient, id.String())
			wss.responseWrite(uuid.Nil(), 3, "add client "+Me.userId.String(), Me.userId, client, uuid.Nil(), 0)
		} else {
			wss.responseWrite(uuid.Nil(), 4, "closed webSocket client "+Me.userId.String(), Me.userId, client, uuid.Nil(), 0)
		}
	}

	if addNewClient {
		fmt.Println("add client "+strings.Join(allClient, " || "))
		wss.responseWrite(uuid.Nil(), 3, "add client "+strings.Join(allClient, " || "), Me.userId, Me, uuid.Nil(), 0)
	}
}
