package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/models"
	bot_utils "rezvin-pro-bot/src/utils/bot"
	utils_context "rezvin-pro-bot/src/utils/context"
)

func (bot *bot) chatIdMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		chatId := bot_utils.GetChatID(update)
		userId := bot_utils.GetUserID(update)

		user := bot.userRepository.GetById(ctx, userId)

		if user != nil && user.ChatId != chatId {
			user.ChatId = chatId
			bot.userRepository.UpdateById(ctx, userId, models.User{
				ChatId: chatId,
			})
		}

		next(utils_context.GetContextWithChatId(ctx, chatId), b, update)
	}
}
