package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/ivykus/gallery/rand"
)

const (
	// MinBytesPerToken is the minimum number of bytes used for
	// each session token
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserId int
	// Token is set only during session creation. When looking up a session
	// TokenHash should be used.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is the number of bytes used for each session token
	// If this value is not set or less than MinBytesPerToken, it will be
	// ignored and set to MinBytesPerToken const
	BytesPerToken int
}

func (ss *SessionService) Create(userId int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := &Session{
		UserId:    userId,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	row := ss.DB.QueryRow(
		`INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE SET token_hash = $2
		RETURNING id`,
		userId,
		session.TokenHash,
	)
	err = row.Scan(&session.ID)

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	token = ss.hash(token)

	user := User{}

	row := ss.DB.QueryRow(
		`SELECT users.id, users.email, users.password_hash
		FROM sessions
		JOIN users
		ON sessions.user_id = users.id
		WHERE sessions.token_hash = $1`,
		token,
	)

	if err := row.Scan(&user.Id, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	token = ss.hash(token)
	_, err := ss.DB.Exec(
		`DELETE FROM sessions
		WHERE token_hash = $1`,
		token,
	)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	hashedBytes := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hashedBytes[:])
}
