package utils_context

import "context"

func GetContextWithChatId(ctx context.Context, chatId int64) context.Context {
	return context.WithValue(ctx, "chatId", chatId)
}

func GetChatIdFromContext(ctx context.Context) int64 {
	result := ctx.Value("chatId")

	if result == nil {
		panic("ChatId not found in context. Error in code")
	}

	return result.(int64)
}
