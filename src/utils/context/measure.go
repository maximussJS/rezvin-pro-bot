package utils_context

import (
	"context"
	"rezvin-pro-bot/src/models"
)

func GetContextWithMeasure(ctx context.Context, program *models.Measure) context.Context {
	return context.WithValue(ctx, "Measure", program)
}

func GetMeasureFromContext(ctx context.Context) *models.Measure {
	result := ctx.Value("Measure")

	if result == nil {
		panic("Measure not found in context. Error in code")
	}

	return result.(*models.Measure)
}
