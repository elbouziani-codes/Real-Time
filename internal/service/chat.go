package service

/* import (
	//	"fmt"
	"context"
	"realTime/internal/domain"
	// "realTime/database"
)

type chatService struct {
	repo RepoChat
}


type RepoChat interface {
	GetChat(context.Context, []string) (domain.ChatRoom, error)
	CreateChat(context.Context, string) (domain.ChatRoom, error)
	GetMessage(context.Context) ([]domain.Message, error)
	AddUsersToChat(context.Context, []domain.User) (error)
}


// []string
func NewService(repo RepoChat) *chatService {
	return &chatService{repo: repo}
}

func (svc *chatService) CreatChat(users []string) (string, error) {
	return "", nil
}

func (svc *chatService) GetChat(users []string)  {
		chatID, err := svc.repo.GetChat(context.Context, []string) (string, error)

		if err != nil  {
				return "", err
		}
		return chatID, nil

} */
