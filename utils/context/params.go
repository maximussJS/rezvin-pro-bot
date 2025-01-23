package utils_context

import (
	"context"
	bot_types "rezvin-pro-bot/types/bot"
)

func GetContextWithParams(ctx context.Context, params *bot_types.Params) context.Context {
	return context.WithValue(ctx, "QueryParams", params)
}

func GetParamsFromContext(ctx context.Context) *bot_types.Params {
	result := ctx.Value("QueryParams")

	if result == nil {
		panic("Params not found in context. Error in code")
	}

	return result.(*bot_types.Params)
}
