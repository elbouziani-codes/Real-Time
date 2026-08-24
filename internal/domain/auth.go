package domain

import (
	"strings"
	"uuid"
)

const (
	EmailType = iota
	UserNameType
)

type Credentials struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt int64
	ExpireAt  int64
}

/// to remove that shit later


func LoginValidation(creds *Credentials) (int, error) {
	creds.Identifier = strings.ToLower(creds.Identifier)	
	idType := UserNameType
	if emailRegex.MatchString(creds.Identifier) {
		idType = EmailType
	} else if !usernameRegex.MatchString(creds.Identifier) {
		return 0, Error{Message: "Invalid Identifier Format", Code: BadFormatCode}
	}

	return idType, nil
}
