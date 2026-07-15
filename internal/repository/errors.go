package repository

import (
	"errors"
	"strings"
	"database/sql"
	"realTime/internal/domain"
	"github.com/mattn/go-sqlite3"
)

func TranslateError(constraintError error) error {
	var err sqlite3.Error
	if errors.Is(sql.ErrNoRows, constraintError)  {
		return  domain.Error{Field: "", Message: "NotFound", Code: domain.NotFoundCode}
	}
	if errors.As(constraintError, &err) {
		message := err.Error()
		switch {
		case strings.Contains(message, "email"):	
				return domain.Error{Field: "email", Message: "already taken", Code: domain.ConflictCode}
		case strings.Contains(message, "nick_name"):
				return domain.Error{Field: "nickname", Message: "already used", Code: domain.ConflictCode}
		default:
			return err
		}
	}
	return err
}
