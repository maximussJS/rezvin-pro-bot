package utils_context

import (
	"context"
	"rezvin-pro-bot/src/models"
)

func GetContextWithUserMeasure(ctx context.Context, program *models.UserMeasure) context.Context {
	return context.WithValue(ctx, "UserMeasure", program)
}

func GetUserMeasureFromContext(ctx context.Context) *models.UserMeasure {
	result := ctx.Value("UserMeasure")

	if result == nil {
		panic("UserMeasure not found in context. Error in code")
	}

	return result.(*models.UserMeasure)
}
