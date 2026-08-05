package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"realTime/internal/domain"
	"strings"

	"github.com/mattn/go-sqlite3"
)

func TranslateError(constraintError error) error {
	fmt.Println(constraintError)
	if constraintError == nil {
		return constraintError
	}
	var err sqlite3.Error
	if sql.ErrNoRows == constraintError {
		return domain.Error{Message: "not found", Code: domain.NotFoundCode}
	}
	if strings.Contains(constraintError.Error(), "parent") {
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
			fmt.Println(err)
			return domain.Error{Message: "unexpected error", Code: domain.UnexpectedCode}
		}
	}
	fmt.Println(constraintError)

	return domain.Error{Message: "unexpected error", Code: domain.UnexpectedCode}

}
