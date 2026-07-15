package domain

import (
	"fmt"
)

const (
	UnexpectedCode  = iota
	NotFoundCode 
	BadFormatCode
	ConflictCode	
	UnauthorizedCode
)

type Error struct {
		Field string // e	
		Message string 
		Code int
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message);
} 
