package bot

import (
	"context"
	"fmt"
)

func (bot *bot) Shutdown(ctx context.Context) error {
	err := bot.server.Shutdown(ctx)

	if err != nil {
		return fmt.Errorf("error while shutting down server: %w", err)
	}

	bot.senderService.SendSafe(ctx, bot.bot, bot.config.AlertChatId(), "Бот вимкнено\\! Схоже сталась критична помилка")

	return nil
}
