package utils_context

import "context"

func GetContextWithLimit(ctx context.Context, limit int) context.Context {
	return context.WithValue(ctx, "Limit", limit)
}

func GetLimitFromContext(ctx context.Context) int {
	result := ctx.Value("Limit")

	if result == nil {
		panic("Limit not found in context. Error in code")
	}

	return result.(int)
}
