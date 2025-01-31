package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	bot_utils "rezvin-pro-bot/src/utils/bot"
	"rezvin-pro-bot/src/utils/messages"
)

func (bot *bot) timeoutMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		timeoutDuration := bot.config.RequestTimeout()
		chatId := bot_utils.GetChatID(update)

		childCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		doneCh := make(chan struct{})

		go func() {
			next(childCtx, b, update)
			close(doneCh)
		}()

		select {
		case <-childCtx.Done():
			bot.senderService.Send(ctx, b, chatId, messages.RequestTimeoutMessage())
			return
		case <-doneCh:
			return
		}
	}
}
