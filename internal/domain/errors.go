package domain

import (
	"fmt"
)

const (
	UnexpectedCode  = iota
	NotFoundCode 
	BadFormatCode
	ConflictCode	
)

type ValidationError struct {
		Field string // e	
		Message string 
		Code int
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message);
} 
