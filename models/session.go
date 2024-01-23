package models

import "database/sql"

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
}

func (ss *SessionService) Create() (*Session, error) {
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}
