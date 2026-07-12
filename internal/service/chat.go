package service

import (
	//	"fmt"
	"realTime/internal/domain"
	// "realTime/database"
)

type chatService struct {
	Repo *domain.RepoChat
}

func NewService(repo *domain.RepoChat) *chatService {
	return &chatService{Repo: repo}
}

func (svc *chatService) CreatChat(users []string) {
}

func (svc *chatService) GetChat(users []string) {
}
