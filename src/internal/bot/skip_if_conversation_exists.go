package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func (bot *bot) skipIfConversationExistsMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		chatId := bot_utils.GetChatID(update)

		if bot.conversationService.IsConversationExists(chatId) {
			return
		}

		next(ctx, b, update)
	}
}
