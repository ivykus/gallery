package models

import (
	"database/sql"
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
		UserId: userId,
		Token:  token,
	}
	return session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}
