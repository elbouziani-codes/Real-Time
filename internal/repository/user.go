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

type scanner interface {
	Scan(...any) error
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

/// MUST BE FIXED LATER
const GetAllUserQuery = `SELECT id, email, nick_name, last_name, first_name FROM users LIMIT  ? OFFSET ?`

func (u *UserRepo) GetUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	users := []domain.User{}
	row, err := u.db.QueryContext(ctx, GetAllUserQuery)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		user, err := scanUser(row)
		if err != nil {
			return nil, sqlite.TranslateError(err)
		}
		users = append(users, user)
	}
	err = row.Err()
	if err != nil {
		return nil, sqlite.TranslateError(err)
	}
	return users, nil
}
