package bot

import (
	"context"
	"fmt"
	"time"
)

func (bot *bot) shutdown() {
	bot.logger.Log("Shutting down gracefully...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	bot.logger.Log("Bot closed successfully")

	err := bot.server.Shutdown(shutdownCtx)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("Failed to shutdown server gracefully: %s", err))
	}

	bot.logger.Log("Server closed successfully")
}
