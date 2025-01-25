package utils_context

import (
	"context"
	"rezvin-pro-bot/src/types"
)

func GetContextWithParams(ctx context.Context, params *types.Params) context.Context {
	return context.WithValue(ctx, "QueryParams", params)
}

func GetParamsFromContext(ctx context.Context) *types.Params {
	result := ctx.Value("QueryParams")

	if result == nil {
		panic("Params not found in context. Error in code")
	}

	return result.(*types.Params)
}
