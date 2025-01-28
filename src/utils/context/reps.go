package utils_context

import (
	"context"
	"rezvin-pro-bot/src/constants"
)

func GetContextWithReps(ctx context.Context, reps constants.Reps) context.Context {
	return context.WithValue(ctx, "Reps", reps)
}

func GetRepsFromContext(ctx context.Context) constants.Reps {
	result := ctx.Value("Reps")

	if result == nil {
		panic("Reps not found in context. Error in code")
	}

	return result.(constants.Reps)
}
