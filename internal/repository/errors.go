package repository

import (
	"fmt"
	"errors"
	"strings"
	"database/sql"
	"realTime/internal/domain"
	"github.com/mattn/go-sqlite3"
)

func TranslateError(constraintError error) error {
	var err sqlite3.Error
	if errors.Is(sql.ErrNoRows, constraintError)  {
		fmt.Println("test")
		return err
	}
	if errors.As(constraintError, &err) {
		message := err.Error()
		switch {
		case strings.Contains(message, "email"):	
				return domain.ValidationError{Field: "email", Message: "already taken", Code: domain.ConflictCode}
		case strings.Contains(message, "nick_name"):
				return domain.ValidationError{Field: "nickname", Message: "already used", Code: domain.ConflictCode}
		default:
			return err
		}
	}
	fmt.Println(err)
	return err
}
