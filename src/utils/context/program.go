package utils_context

import (
	"context"
	"rezvin-pro-bot/src/models"
)

func GetContextWithProgram(ctx context.Context, program *models.Program) context.Context {
	return context.WithValue(ctx, "Program", program)
}

func GetProgramFromContext(ctx context.Context) *models.Program {
	result := ctx.Value("Program")

	if result == nil {
		panic("Program not found in context. Error in code")
	}

	return result.(*models.Program)
}
