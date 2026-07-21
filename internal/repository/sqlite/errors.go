package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"realTime/internal/domain"
	"strings"
)

func TranslateError(constraintError error) error {
	fmt.Println(constraintError)
	var err sqlite3.Error
	if sql.ErrNoRows == constraintError {
		return domain.Error{Message: "not found", Code: domain.NotFoundCode}
	}
	if errors.As(constraintError, &err) {
		message := err.Error()
		switch {
		case strings.Contains(message, "email"):
			return domain.Error{Message: "email lready taken", Code: domain.ConflictCode}
		case strings.Contains(message, "nick_name"):
			return domain.Error{Message: "nickname already used", Code: domain.ConflictCode}
		default:
			return domain.Error{Message: "unexpected error", Code: domain.UnexpectedCode}
		}
	}
	return domain.Error{Message: "unexpected error", Code: domain.UnexpectedCode}

}
