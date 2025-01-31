package bot

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"rezvin-pro-bot/src/constants"
	"time"
)

func (bot *bot) Start() {
	defer bot.shutdownWaitGroup.Done()

	bot.shutdownWaitGroup.Add(1)

	go func() {
		if bot.config.AppEnv() == constants.DevelopmentEnv {
			bot.startPolling()
		} else {
			bot.startWebhook()
		}
	}()

	for _, chatId := range bot.config.AdminChatIds() {
		bot.senderService.Send(bot.shutdownContext, bot.bot, chatId, "Бот запустився і готовий до роботи\\!")
	}

	for {
		select {
		case <-bot.shutdownContext.Done():
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			for _, chatId := range bot.config.AdminChatIds() {
				bot.senderService.SendSafe(ctx, bot.bot, chatId, "Бот вимкнено\\! Схоже сталась критична помилка")
			}
			bot.shutdown()
			cancel()
			return
		}
	}
}

func (bot *bot) startWebhook() {
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

	bot.bot.StartWebhook(bot.shutdownContext)
}

func (bot *bot) startPolling() {
	bot.logger.Log("Bot started in polling mode")

	bot.bot.Start(bot.shutdownContext)
}
