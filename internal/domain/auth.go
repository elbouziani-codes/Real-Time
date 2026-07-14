package domain

import (
	"regexp"
)



type Credentials struct {
	EmailOrNickName string `json:"identifier"`
	Password        string `json:"password"`
}

type Session struct {	
	ID string 
	UserID string	
	CreatedAt int64
	ExpireAt int64
}


/// to remove that shit later
var (
	usernameRegex = regexp.MustCompile(`^[\p{L}\p{N}_]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)


