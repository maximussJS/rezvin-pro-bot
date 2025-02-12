package utils_context

import (
	"context"
	"rezvin-pro-bot/src/models"
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

func GetContextWithCurrentUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, "CurrentUser", user)
}

func GetCurrentUserFromContext(ctx context.Context) *models.User {
	result := ctx.Value("CurrentUser")

	if result == nil {
		panic("CurrentUser not found in context. Error in code")
	}

	return result.(*models.User)
}
