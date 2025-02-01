package bot

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"rezvin-pro-bot/src/constants"
	"time"
)

func (bot *bot) Start(ctx context.Context) {
	bot.senderService.Send(ctx, bot.bot, bot.config.AlertChatId(), "Бот запустився і готовий до роботи\\!")

	if bot.config.AppEnv() == constants.DevelopmentEnv {
		bot.startPolling(ctx)
	} else {
		bot.startWebhook(ctx)
	}

}

func (bot *bot) startWebhook(ctx context.Context) {
	tlsConfig := &tls.Config{
		ClientAuth: tls.NoClientCert,
		MinVersion: tls.VersionTLS11,
	}

	port := bot.config.HttpPort()

	bot.server = http.Server{
		Addr:         port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      bot.bot.WebhookHandler(),
	}

	go func() {
		bot.logger.Log(fmt.Sprintf("Starting https server on port %s", port))
		if err := bot.server.ListenAndServeTLS(bot.config.SSLCertPath(), bot.config.SSLKeyPath()); err != nil && err != http.ErrServerClosed {
			bot.logger.Fatal(fmt.Sprintf("Failed to start https server: %s", err))
		}
	}()

	bot.logger.Log("Bot started in webhook mode")

	bot.bot.StartWebhook(ctx)
}

func (bot *bot) startPolling(ctx context.Context) {
	bot.logger.Log("Bot started in polling mode")

	bot.bot.Start(ctx)
}
