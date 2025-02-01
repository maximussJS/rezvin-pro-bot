package utils_context

import (
	"context"
	"rezvin-pro-bot/src/models"
)

func GetContextWithUserResult(ctx context.Context, userExercise *models.UserResult) context.Context {
	return context.WithValue(ctx, "UserResult", userExercise)
}

func GetUserResultFromContext(ctx context.Context) *models.UserResult {
	result := ctx.Value("UserResult")

	if result == nil {
		panic("UserResult not found in context. Error in code")
	}

	return result.(*models.UserResult)
}
