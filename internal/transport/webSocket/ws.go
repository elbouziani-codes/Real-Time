package ws

import (
	"context"
	"errors"
	"fmt"
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
	clients map[crypto.UUID]*Client
}
type UserService interface {
	GetByID(context.Context, crypto.UUID) (domain.User, error)
}

type ChatService interface {
	IsChatNotFound(error) bool
	CheckRoomChat(context.Context, []crypto.UUID) (crypto.UUID, error)
	CreateRoomChat(context.Context, []crypto.UUID) (crypto.UUID, error)
	SendMessageRoomChat(context.Context, string, crypto.UUID, crypto.UUID) (crypto.UUID, int64, error)
}

func NewHandleWs(svcChat ChatService, svcUser UserService) *HandlerWs {
	return &HandlerWs{
		svcChat: svcChat,
		svcUser: svcUser,
		hub: &hub{
			clients: make(map[crypto.UUID]*Client),
		},
	}
}



type Client struct {
	userId crypto.UUID
	conn   *websocket.Conn
	mu     sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wss *HandlerWs) ChatWs(w http.ResponseWriter, r *http.Request) {
	client, err := wss.AddClient(w, r)
	if err != nil {
		return
	}
	go wss.engineMessages(r.Context(), client)
}

func (wss *HandlerWs) engineMessages(ctx context.Context, sender *Client) {
	defer func() {
		wss.DeleteClient(sender)
	}()
	for {
		err, wsRequest := wss.readInputMessage(ctx, sender)
		if err != nil && err.Error() == "Close" {
			return
		} else if err != nil {
			continue
		}
		chatId, err := wss.GetRoom(ctx, sender.userId,wsRequest.Destination)
		if err != nil {
			wss.responseWrite(crypto.Nil, 404, err.Error(), sender.userId, sender , chatId ,0)
			continue
		}
		
		if wsRequest.RequestType == "typing" {
			fmt.Println("dddddfddddd")
			wss.sendTypingReceiver(wsRequest.Destination, sender.userId , chatId)
			continue
		}
		

		messageID, createAt, err := wss.svcChat.SendMessageRoomChat(ctx, wsRequest.Content, sender.userId, chatId)
		if err != nil {
			wss.responseWrite(crypto.Nil, 404, err.Error(), sender.userId, sender , chatId ,createAt)
			continue
		}

		wss.sendReceiver(wsRequest, messageID, sender.userId ,chatId ,createAt)
	}
}

func (wss *HandlerWs) sendReceiver(wsRequest *domain.WsParsedRequest, messageID crypto.UUID, senderID crypto.UUID , chatId crypto.UUID , create_at int64) {
	wss.hub.mu.RLock()
	receiver := wss.hub.clients[wsRequest.Destination]
	wss.hub.mu.RUnlock()
	if receiver != nil {
		wss.responseWrite(messageID, 200, wsRequest.Content, senderID, receiver,chatId,create_at)
	}
}

func (wss *HandlerWs) sendTypingReceiver(receiverID crypto.UUID, senderID crypto.UUID,chatId crypto.UUID) {
	wss.hub.mu.RLock()
	receiver := wss.hub.clients[receiverID]
	wss.hub.mu.RUnlock()
	if receiver != nil {
		wss.responseWrite(crypto.Nil, 2, "typing", senderID, receiver , chatId ,0)
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
		wss.responseWrite(crypto.Nil, 404, err.Error(), client.userId, client , crypto.Nil,0)
		return err, &domain.WsParsedRequest{}
	}
	_, err = wss.svcUser.GetByID(ctx, wsRequestParse.Destination)

	if err != nil {
		wss.responseWrite(crypto.Nil, 404, err.Error(), client.userId, client, crypto.Nil,0)
		return err, &domain.WsParsedRequest{}
	}

	if client.userId.Value == wsRequestParse.Destination.Value {
		wss.responseWrite(crypto.Nil, 404, "Error in sender == Destination", client.userId, client, crypto.Nil,0)
		return domain.Error{Message: "Error in sender == Destination", Code: domain.ConflictCode}, &domain.WsParsedRequest{}
	}
fmt.Println("conetnt", wsRequestParse.Content)
	return nil, wsRequestParse
}

func (wss *HandlerWs) responseWrite(id crypto.UUID, codeError int, content string, SenderID crypto.UUID, receiver *Client, chatID crypto.UUID, createdAt int64) {
	MessageOutput := domain.MessageOutput{ID: id, Code: codeError, Content: content, Sender: SenderID, ChatID: chatID, Created_at: createdAt}
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	err := receiver.conn.WriteJSON(&MessageOutput)
	if err != nil {
		go func() {
			wss.DeleteClient(receiver)
		}()
	}
}

func (wss *HandlerWs) GetRoom(ctx context.Context, userId crypto.UUID, ReceiverID crypto.UUID) (crypto.UUID, error) {
	userIDs := []crypto.UUID{userId, ReceiverID}
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

func (wss *HandlerWs) DeleteClient(client *Client) {
	wss.hub.mu.Lock()

	current, exists := wss.hub.clients[client.userId]

	if !exists || current != client {
		wss.hub.mu.Unlock()
		return
	}

	delete(wss.hub.clients, client.userId)

	wss.hub.mu.Unlock()

	client.conn.Close()
}

func (wss *HandlerWs) AddClient(w http.ResponseWriter, r *http.Request) (*Client, error) {
	userIDAny := r.Context().Value("user_id")
	userID, _ := userIDAny.(crypto.UUID)
	wss.hub.mu.Lock()
	oldClient, ok := wss.hub.clients[userID]
	wss.hub.mu.Unlock()
	if ok {
		wss.responseWrite(crypto.Nil, 1, "Close WebSocket Connection and Disable Chat Page", oldClient.userId, oldClient ,crypto.Nil,0)
		wss.DeleteClient(oldClient)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	wss.hub.mu.Lock()
	defer wss.hub.mu.Unlock()

	sender := &Client{
		userId: userID,
		conn:   conn,
	}

	wss.hub.clients[userID] = sender
	return sender, nil
}
