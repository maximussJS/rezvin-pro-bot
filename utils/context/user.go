package utils_context

import (
	"context"
	"rezvin-pro-bot/models"
)

func GetContextWithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, "User", user)
}

func GetUserFromContext(ctx context.Context) *models.User {
	result := ctx.Value("User")

	if result == nil {
		panic("User not found in context. Error in code")
	}

	return result.(*models.User)
}
