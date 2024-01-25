package context

import (
	"context"

	"github.com/ivykus/gallery/models"
)

type key int

const (
	userKey key = iota
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	userValue := ctx.Value(userKey)
	user, ok := userValue.(*models.User)
	if !ok {
		// if nothing stored in the conext, or value is invalid
		return nil
	}
	return user
}
