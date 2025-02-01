package main

import (
	"context"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/internal/bot"
	"rezvin-pro-bot/src/internal/db"
	"rezvin-pro-bot/src/services"
	"rezvin-pro-bot/src/types"
	"rezvin-pro-bot/src/utils"
)

type runAppDependencies struct {
	dig.In

	ShutdownContext context.Context `name:"ShutdownContext"`

	Database db.IDatabase `name:"Database"`

	LockService         services.ILockService         `name:"LockService"`
	ShutdownService     services.IShutdownService     `name:"ShutdownService"`
	ConversationService services.IConversationService `name:"ConversationService"`

	Bot bot.IBot `name:"Bot"`
}

func StartApplication(container *dig.Container) {
	err := container.Invoke(func(deps runAppDependencies) {
		botShutdownCallback := types.NewShutdownCallback(
			"Bot",
			func(ctx context.Context) error {
				return deps.Bot.Shutdown(ctx)
			},
			4,
		)

		conversationServiceShutdownCallback := types.NewShutdownCallback(
			"ConversationService",
			func(ctx context.Context) error {
				return deps.ConversationService.Shutdown(ctx)
			},
			2,
		)

		lockServiceShutdownCallback := types.NewShutdownCallback(
			"LockService",
			func(ctx context.Context) error {
				return deps.LockService.Shutdown(ctx)
			},
			1,
		)

		databaseShutdownCallback := types.NewShutdownCallback(
			"Database",
			func(ctx context.Context) error {
				return deps.Database.Shutdown(ctx)
			},
			3,
		)

		deps.ShutdownService.AddShutdownCallback(botShutdownCallback)
		deps.ShutdownService.AddShutdownCallback(conversationServiceShutdownCallback)
		deps.ShutdownService.AddShutdownCallback(lockServiceShutdownCallback)
		deps.ShutdownService.AddShutdownCallback(databaseShutdownCallback)

		go deps.Bot.Start(deps.ShutdownContext)
	})

	utils.PanicIfError(err)
}
