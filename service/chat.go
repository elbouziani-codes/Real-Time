package service

import (
	//	"fmt"
	"realTime/domain"
	"sync"
	// "realTime/database"
)

type chatService struct {
	Repo *domain.RepoChat
}

func NewService(repo *domain.RepoChat) *chatService {
	return &chatService{Repo: repo}
}


func (svc *chatService ) CreatChat(users [][]byte) {
	
}

func (svc *chatService ) GetChat(users [][]byte) {

}




func (svc *chatService ) chatLoop() {

}
