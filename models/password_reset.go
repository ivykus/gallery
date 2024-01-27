package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/ivykus/gallery/rand"
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
	email = strings.ToLower(email)
	var userID int

	row := ps.DB.QueryRow(`
		SELECT id
		FROM users
		WHERE email = $1
	`, email)
	err := row.Scan(&userID)
	if err != nil {
		// TODO consider handling different types of errors for when the user
		// is not found and if there is db error
		return nil, fmt.Errorf("create: %w", err)
	}

	bytesPerToken := ps.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	duration := ps.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserId:    userID,
		Token:     token,
		TokenHash: ps.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	row = ps.DB.QueryRow(`
		INSERT INTO passwords_reset (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO
		UPDATE SET token_hash = $2, expires_at = $3
		RETURNING id
	`, pwReset.UserId, pwReset.TokenHash, pwReset.ExpiresAt)

	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwReset, nil
}

func (ps *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Consume")
}

func (ps *PasswordResetService) hash(token string) string {
	hashedBytes := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hashedBytes[:])
}
