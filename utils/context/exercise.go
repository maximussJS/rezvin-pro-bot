package utils_context

import (
	"context"
	"rezvin-pro-bot/models"
)

func GetContextWithExercise(ctx context.Context, exercise *models.Exercise) context.Context {
	return context.WithValue(ctx, "Exercise", exercise)
}

func GetExerciseFromContext(ctx context.Context) *models.Exercise {
	result := ctx.Value("Exercise")

	if result == nil {
		panic("Exercise not found in context. Error in code")
	}

	return result.(*models.Exercise)
}
