package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	. "rezvin-pro-bot/utils"
)

func MustSendMessage(ctx context.Context, b *tg_bot.Bot, params *tg_bot.SendMessageParams) *models.Message {
	msg, err := b.SendMessage(ctx, params)

	PanicIfError(err)

	return msg
}

func MustAnswerCallbackQuery(ctx context.Context, b *tg_bot.Bot, update *models.Update) bool {
	result, err := b.AnswerCallbackQuery(ctx, &tg_bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	PanicIfError(err)

	return result
}
