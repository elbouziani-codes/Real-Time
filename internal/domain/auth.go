package domain

import (
	"realTime/crypto"
	"regexp"
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
	ID        crypto.UUID
	UserID    crypto.UUID
	CreatedAt int64
	ExpireAt  int64
}

/// to remove that shit later
var (
	usernameRegex = regexp.MustCompile(`^[\p{L}\p{N}_]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func LoginValueidation(creds Credentials) (int, error) {
	idType := UserNameType
	if emailRegex.MatchString(creds.Identifier) {
		idType = EmailType
	} else if !usernameRegex.MatchString(creds.Identifier) {
		return 0, Error{Message: "Invalid Identifier Format", Code: BadFormatCode}
	}

	return idType, nil
}
