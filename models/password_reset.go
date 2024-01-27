package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserId    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	// Duration is the length of time a password reset token is valid for
	Duration time.Duration
}

func (ps *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Create")
}

func (ps *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Consume")
}
