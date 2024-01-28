package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

var (
	ErrEmailTaken = fmt.Errorf("models: email is already in use")
)

func (us *UserService) CreateUser(email, password string) (*User, error) {
	email = strings.ToLower(email)

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating password hash: %w", err)
	}
	user := User{
		Email:        email,
		PasswordHash: string(hashBytes),
	}

	raw := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id;	
	`, user.Email, user.PasswordHash)

	err = raw.Scan(&user.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("inserting user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	raw := us.DB.QueryRow(`
		SELECT id, password_hash
		FROM users
		WHERE email = $1
	`, email)

	err := raw.Scan(&user.Id, &user.PasswordHash)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("authenticate: %w", err)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return &user, nil
}

func (us *UserService) UpdatePassword(userID int, password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashBytes)
	_, err = us.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1
	`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}
