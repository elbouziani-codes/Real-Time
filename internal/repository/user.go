package repository

import (
	"context"
	"database/sql"
	"errors"

	"realTime/crypto"
	"realTime/internal/domain"
	"realTime/internal/repository/sqlite"
)

type UserRepo struct {
	db DBTX
}



func NewUserRepo(db DBTX) *UserRepo {
	return &UserRepo{db: db}
}

func scanUser(row scanner) (domain.User, error) {
	user := domain.User{}
	err := row.Scan(
		&user.ID.Value,
		&user.Email,
		&user.Password,
		&user.NickName,
		&user.LastName,
		&user.FirstName,
		&user.Age,
		&user.Gender,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return domain.User{}, sqlite.TranslateError(err)
	}
	return user, nil
}

func scanUserProfile(row scanner) (*domain.UserProfile, error) {
	user := domain.UserProfile{}
	err := row.Scan(
		&user.ID.Value,
		&user.Email,
		&user.NickName,
		&user.LastName,
		&user.FirstName,
		&user.Age,
		&user.Gender)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return &user, nil
}

const CreateUserQuery = `INSERT INTO users 
(id, email, password_hash, nick_name, last_name, first_name, age, gender) VALUES 
(?, ?, ?, ?, ?, ?, ?, ?) `

func (u *UserRepo) CreateUser(ctx context.Context, user domain.User) error {
	_, err := u.db.ExecContext(ctx, CreateUserQuery,
		user.ID.Value,
		user.Email,
		user.Password,
		user.NickName,
		user.LastName,
		user.FirstName,
		user.Age,
		user.Gender)
	if err != nil {
		return sqlite.TranslateError(err)
	}
	return nil
}

const GetByIdQuery = `SELECT id, email, password_hash, nick_name, last_name, first_name, age, gender, created_at, updated_at
FROM users WHERE id = ?`

func (u *UserRepo) GetByID(ctx context.Context, userID crypto.UUID) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, GetByIdQuery, userID.Value)
	return scanUser(row)
}

const GetProfileQuery = `SELECT id, email, nick_name, last_name, first_name, age, gender 
FROM users WHERE id = ?`

func (u *UserRepo) GetUserProfile(ctx context.Context, requesterID, userID crypto.UUID) (*domain.UserProfile, error) {
	row := u.db.QueryRowContext(ctx, GetProfileQuery, userID.Value)
	return scanUserProfile(row)
}

const GetByEmailQuery = `SELECT id, email, password_hash, nick_name, last_name, first_name, age, gender, created_at, updated_at
FROM users WHERE email = ?`

func (u *UserRepo) GetByEmail(ctx context.Context, userEmail string) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, GetByEmailQuery, userEmail)
	return scanUser(row)
}

const GetByNickNameQuery = `SELECT id, email, password_hash, nick_name, last_name, first_name, age, gender, created_at, updated_at
FROM users WHERE nick_name = ?`

// Idont think we will need that
func (u *UserRepo) GetByNickName(ctx context.Context, NickName string) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, GetByNickNameQuery, NickName)
	return scanUser(row)
}

// lastMessageAtColumn computes the newest message shared between the requester
// and one user. It is a correlated subquery rather than a join so a
// conversation holding many messages cannot multiply a user's row, and it reads
// only conversations the requester is a participant of — joining messages on the
// other user's conversation alone would expose activity from chats the requester
// is not in.
const lastMessageAtColumn = `
COALESCE((
	SELECT MAX(M.created_at)
	FROM messages M
	JOIN conversation_participants MINE
	ON MINE.conversation_id = M.conversation_id AND MINE.user_id = ?
	JOIN conversation_participants THEIRS
	ON THEIRS.conversation_id = M.conversation_id AND THEIRS.user_id = U.id
), 0)`

// getAllUsersQueryHead and getAllUsersQueryTail sandwich the optional cursor
// condition. The listing is wrapped in a derived table so the cursor condition
// can reference the last_message_at alias, which a bare WHERE clause cannot.
// Users never messaged sort last on 0, then alphabetically so the tail of the
// list is stable rather than arbitrary.
const getAllUsersQueryHead = `
SELECT id, email, nick_name, last_name, first_name, age, gender, last_message_at
FROM (
	SELECT U.id, U.email, U.nick_name, U.last_name, U.first_name, U.age, U.gender,
	` + lastMessageAtColumn + ` AS last_message_at
	FROM users U
	WHERE U.id != ?
)
`

const getAllUsersQueryTail = `
ORDER BY last_message_at DESC, nick_name ASC
LIMIT ?
`

const getUserCursorQuery = `
SELECT ` + lastMessageAtColumn + `, U.nick_name
FROM users U
WHERE U.id = ?
`

// userCursorKey is the sort position of the user a client last received. The
// nick_name belongs to it because last_message_at alone does not identify a
// single row.
type userCursorKey struct {
	lastMessageAt int
	nickName      string
}

// resolveUserCursor reads the sort position of the cursor user. A cursor naming
// a user that has since been deleted is reported rather than silently returning
// an empty page, which a client could not tell apart from the end of the list.
func (u *UserRepo) resolveUserCursor(ctx context.Context, requesterID, cursor crypto.UUID) (*userCursorKey, error) {
	if cursor == crypto.Nil {
		return nil, nil
	}

	key := userCursorKey{}
	err := u.db.QueryRowContext(ctx, getUserCursorQuery, requesterID.Value, cursor.Value).Scan(&key.lastMessageAt, &key.nickName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.Error{Message: "unknown cursor", Code: domain.NotFoundCode}
	}
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return &key, nil
}

func scanUserContact(row scanner) (domain.UserContact, error) {
	contact := domain.UserContact{}
	err := row.Scan(
		&contact.ID.Value,
		&contact.Email,
		&contact.NickName,
		&contact.LastName,
		&contact.FirstName,
		&contact.Age,
		&contact.Gender,
		&contact.LastMessageAt)
	if err != nil {
		return domain.UserContact{}, sqlite.TranslateError(err)
	}
	return contact, nil
}

// GetUsers returns the people list for userID, ordered by the most recent
// conversation. The requester's own id is bound twice: once to find the shared
// conversations and once to leave themselves out of their own list.
func (u *UserRepo) GetUsers(ctx context.Context, userID crypto.UUID, limit int, cursor crypto.UUID) ([]domain.UserContact, error) {
	cursorKey, err := u.resolveUserCursor(ctx, userID, cursor)
	if err != nil {
		return nil, err
	}

	// Keyset paging: keep the rows that fall strictly after the cursor under the
	// ORDER BY above. nick_name is unique, so it is a total tie-breaker and the
	// cursor row itself is never repeated.
	where := ""
	args := []any{userID.Value, userID.Value}
	if cursorKey != nil {
		where = "WHERE last_message_at < ? OR (last_message_at = ? AND nick_name > ?)"
		args = append(args, cursorKey.lastMessageAt, cursorKey.lastMessageAt, cursorKey.nickName)
	}
	args = append(args, limit)

	users := []domain.UserContact{}
	rows, err := u.db.QueryContext(ctx, getAllUsersQueryHead+where+getAllUsersQueryTail, args...)
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	defer rows.Close()
	for rows.Next() {
		user, err := scanUserContact(rows)
		if err != nil {
			return nil, sqlite.TranslateError(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return users, nil
}
