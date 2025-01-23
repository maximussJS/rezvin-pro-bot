package utils_context

import "context"

func GetContextWithOffset(ctx context.Context, offset int) context.Context {
	return context.WithValue(ctx, "Offset", offset)
}

func GetOffsetFromContext(ctx context.Context) int {
	result := ctx.Value("Offset")

	if result == nil {
		panic("Offset not found in context. Error in code")
	}

	return result.(int)
}
