package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	bot_utils "rezvin-pro-bot/src/utils/bot"
	"rezvin-pro-bot/src/utils/messages"
	"runtime"
)

func (bot *bot) panicRecoveryMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		defer func() {
			if err := recover(); err != nil {
				chatID := bot_utils.GetChatID(update)
				stackSize := bot.config.ErrorStackTraceSizeInKb() * 1024
				stack := make([]byte, stackSize)
				length := runtime.Stack(stack, true)
				stack = stack[:length]

				if ctx.Err() != nil {
					return
				}

				bot.logger.Error(fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack))

				b.SendMessage(ctx, &tg_bot.SendMessageParams{
					ChatID:    chatID,
					Text:      messages.ErrorMessage(),
					ParseMode: tg_models.ParseModeMarkdown,
				})
			}
		}()

		next(ctx, b, update)
	}
}
