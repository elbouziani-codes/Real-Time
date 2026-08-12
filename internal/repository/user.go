package repository

import (
	"context"
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

// GetAllUserQuery lists everyone except the requester, most recently talked with
// first. last_message_at is a correlated subquery rather than a join so a
// conversation holding many messages cannot multiply a user's row, and it reads
// only conversations the requester is a participant of — joining messages on the
// other user's conversation alone would expose activity from chats the requester
// is not in. Users never messaged sort last on 0, then alphabetically so the
// tail of the list is stable rather than arbitrary.
const GetAllUserQuery = `
SELECT U.id, U.email, U.nick_name, U.last_name, U.first_name, U.age, U.gender,
COALESCE((
	SELECT MAX(M.created_at)
	FROM messages M
	JOIN conversation_participants MINE
	ON MINE.conversation_id = M.conversation_id AND MINE.user_id = ?
	JOIN conversation_participants THEIRS
	ON THEIRS.conversation_id = M.conversation_id AND THEIRS.user_id = U.id
), 0) AS last_message_at
FROM users U
WHERE U.id != ?
ORDER BY last_message_at DESC, U.nick_name ASC
LIMIT ? OFFSET ?
`

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
func (u *UserRepo) GetUsers(ctx context.Context, userID crypto.UUID, limit, offset int) ([]domain.UserContact, error) {
	users := []domain.UserContact{}
	rows, err := u.db.QueryContext(ctx, GetAllUserQuery, userID.Value, userID.Value, limit, offset)
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
