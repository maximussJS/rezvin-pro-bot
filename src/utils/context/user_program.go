package utils_context

import (
	"context"
	"rezvin-pro-bot/src/models"
)

func GetContextWithUserProgram(ctx context.Context, userProgram *models.UserProgram) context.Context {
	return context.WithValue(ctx, "UserProgram", userProgram)
}

func GetUserProgramFromContext(ctx context.Context) *models.UserProgram {
	result := ctx.Value("UserProgram")

	if result == nil {
		panic("UserProgram not found in context. Error in code")
	}

	return result.(*models.UserProgram)
}
