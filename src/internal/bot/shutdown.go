package bot

import (
	"context"
	"fmt"
	"rezvin-pro-bot/src/globals"
)

func (bot *bot) Shutdown(ctx context.Context) error {
	err := bot.server.Shutdown(ctx)

	if err != nil {
		return fmt.Errorf("error while shutting down server: %w", err)
	}

	bot.senderService.SendSafe(ctx, bot.bot, bot.config.AlertChatId(), fmt.Sprintf("Бот %s вимкнено\\! Схоже сталась критична помилка", globals.AdminName))

	return nil
}
