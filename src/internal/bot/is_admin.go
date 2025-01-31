package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	utils_context "rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/messages"
)

func (bot *bot) isAdminMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		user := utils_context.GetCurrentUserFromContext(ctx)

		if user.IsNotAdmin() {
			chatId := utils_context.GetChatIdFromContext(ctx)

			bot.senderService.Send(ctx, b, chatId, messages.AdminOnlyMessage())
			return
		}

		next(ctx, b, update)
	}
}
