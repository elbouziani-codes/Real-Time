package domain

import (
	"realTime/crypto"
	"regexp"
	"strconv"
	"strings"
)

type User struct {
	ID        crypto.UUID
	NickName  string
	LastName  string
	FirstName string
	Email     string
	Password  string
	Gender    string
	Age       int
	CreatedAt int
	UpdatedAt int
}

type UserProfile struct {
	ID        crypto.UUID
	Email  	string
	NickName  string
	LastName  string
	FirstName string
	Gender    string
	Age       int
	CreatedAt int
}

// UserContact is a user as they appear in the people list. LastMessageAt holds
// the newest message in a conversation shared with the requester and is 0 for
// someone they have never messaged; it is the key the list is ordered on, so it
// travels with the profile rather than staying hidden in the query.
type UserContact struct {
	UserProfile
	LastMessageAt int
}

// DefaultUserListLimit and MaxUserListLimit bound the people list. A caller may
// ask for a smaller page but cannot ask for an unbounded one.
const (
	DefaultUserListLimit = 50
	MaxUserListLimit     = 100
)

// ValidateUserListPaging clamps the paging parameters. An unparsable or negative
// limit falls back to the default instead of erroring, so a hand-edited link
// still renders a list rather than a 400.
func ValidateUserListPaging(rawLimit, rawCursor string) (int, crypto.UUID, error) {
	limit := DefaultUserListLimit
	if parsed, err := strconv.Atoi(strings.TrimSpace(rawLimit)); err == nil && parsed > 0 {
		limit = parsed
	}
	if limit > MaxUserListLimit {
		limit = MaxUserListLimit
	}

	cursor, err := ValidateUserCursor(rawCursor)
	if err != nil {
		return 0, crypto.Nil, err
	}

	return limit, cursor, nil
}

// ValidateUserCursor parses the id of the last user the client already holds.
// Paging on that user's sort position rather than a row count keeps a page from
// skipping or repeating users when the ordering shifts mid-scroll. An empty
// value asks for the first page, so an opening request carries no cursor.
func ValidateUserCursor(rawCursor string) (crypto.UUID, error) {
	rawCursor = strings.TrimSpace(rawCursor)
	if rawCursor == "" {
		return crypto.Nil, nil
	}

	cursor, err := crypto.ParseUUID(rawCursor)
	if err != nil {
		return crypto.Nil, Error{Message: "invalid cursor", Code: BadFormatCode}
	}
	return cursor, nil
}
type RegisterRequest struct {
	NickName  string `json:"nick_name"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
}

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]*$`) // would deleted later

func ValueidateUserInfo(registerRequest RegisterRequest) (User, error) {
	var user User
	user.Email = strings.TrimSpace(registerRequest.Email)
	user.Password = registerRequest.Password
	user.NickName = strings.TrimSpace(registerRequest.NickName)
	user.LastName = strings.TrimSpace(registerRequest.LastName)
	user.FirstName = strings.TrimSpace(registerRequest.FirstName)
	user.Gender = strings.TrimSpace(registerRequest.Gender)
	user.Age = registerRequest.Age
	if len(user.Password) < 8 || len(user.Password) > 20 {
		return user, Error{Message: "password length must be between 8 and 20", Code: BadFormatCode}
	}

	if len(user.Email) < 6 || len(user.Email) > 75 {
		return user, Error{Message: "email length must be between 6 and 75", Code: BadFormatCode} // 400
	}

	if !emailRegex.MatchString(user.Email) {
		return user, Error{Message: "email invalid format", Code: BadFormatCode}
	}

	if !usernameRegex.MatchString(user.NickName) {
		return user, Error{Message: "nickname invalid format", Code: BadFormatCode}
	}

	if len(user.NickName) > 20 || len(user.NickName) < 2 {
		return user, Error{Message: "nickname length must be between 2 and 20", Code: BadFormatCode}
	}

	if !nameRegex.MatchString(user.FirstName) || !nameRegex.MatchString(user.LastName) {
		return user, Error{Message: "invalid FirstName or LastName", Code: BadFormatCode}
	}

	if (len(user.FirstName) > 25 || len(user.FirstName) < 2) || (len(user.LastName) > 25 || len(user.LastName) < 2) {
		return user, Error{Message: "firstName and LastName length must be between 2 and 25", Code: BadFormatCode}
	}
	if user.Gender != "man" && user.Gender != "woman" {
		return user, Error{Message: "value gender is not valid", Code: BadFormatCode}
	}
	if user.Age < 14 {
		return user, Error{Message: "you are too old", Code: BadFormatCode}
	}
	if user.Age > 200 {
		return user, Error{Message: "invalid age", Code: BadFormatCode}
	}
	return user, nil
}
