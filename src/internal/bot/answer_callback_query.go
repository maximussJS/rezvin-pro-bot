package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	utils_context "rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
)

func (bot *bot) answerCallbackQueryMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		answerResult := bot.senderService.AnswerCallbackQuery(ctx, b, update)

		if !answerResult {
			chatId := utils_context.GetChatIdFromContext(ctx)

			bot.logger.Error(fmt.Sprintf("Failed to answer callback query: %s", update.CallbackQuery.ID))

			msg := messages.ErrorMessage()
			kb := inline_keyboards.StartOk()

			bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
			return
		}

		next(ctx, b, update)
	}
}
