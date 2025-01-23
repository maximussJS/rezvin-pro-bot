package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	. "rezvin-pro-bot/utils"
)

func SendMessage(ctx context.Context, b *tg_bot.Bot, chatId int64, message string) *models.Message {
	msg, err := b.SendMessage(ctx, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      message,
		ParseMode: models.ParseModeMarkdown,
	})

	PanicIfError(err)

	return msg
}

func SendMessageWithInlineKeyboard(ctx context.Context, b *tg_bot.Bot, chatId int64, message string, kb *models.InlineKeyboardMarkup) *models.Message {
	msg, err := b.SendMessage(ctx, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        message,
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeMarkdown,
	})

	PanicIfError(err)

	return msg
}

func AnswerCallbackQuery(ctx context.Context, b *tg_bot.Bot, update *models.Update) bool {
	result, err := b.AnswerCallbackQuery(ctx, &tg_bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	PanicIfError(err)

	return result
}
