package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
)

func (bot *bot) skipOtherTypesMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		if update.Message == nil && update.CallbackQuery == nil {
			bot.logger.Log(fmt.Sprintf("skipOtherTypesMiddleware: update is not message or callback query. Update: %v", update))
			return
		}

		next(ctx, b, update)
	}
}
