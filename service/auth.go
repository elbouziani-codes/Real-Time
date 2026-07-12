package service

import (
	"realTime/domain"
)

type AuthService struct {
	userRepo domain.UserRepo
}

func (a AuthService) login(credentials domain.Credentials) {
	
}
