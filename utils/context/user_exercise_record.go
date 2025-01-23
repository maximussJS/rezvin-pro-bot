package utils_context

import (
	"context"
	"rezvin-pro-bot/models"
)

func GetContextWithUserExerciseRecord(ctx context.Context, userExercise *models.UserExerciseRecord) context.Context {
	return context.WithValue(ctx, "UserExerciseRecord", userExercise)
}

func GetUserExerciseRecordFromContext(ctx context.Context) *models.UserExerciseRecord {
	result := ctx.Value("UserExerciseRecord")

	if result == nil {
		panic("UserExerciseRecord not found in context. Error in code")
	}

	return result.(*models.UserExerciseRecord)
}
