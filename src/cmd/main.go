package main

import (
	"context"
	"go.uber.org/dig"
	"os"
	"os/signal"
	di2 "rezvin-pro-bot/src/di"
	"rezvin-pro-bot/src/internal/bot"
	"rezvin-pro-bot/src/utils"
	"sync"
	"syscall"
)

type runAppDependencies struct {
	dig.In

	Bot bot.IBot `name:"Bot"`
}

func start(container *dig.Container) {
	err := container.Invoke(func(deps runAppDependencies) {
		go deps.Bot.Start()
	})

	utils.PanicIfError(err)
}

func main() {
	shutdownContext, cancel := context.WithCancel(context.Background())

	defer cancel()

	var wg sync.WaitGroup

	container := di2.BuildContainer()

	container = di2.AppendDependenciesToContainer(container, []di2.Dependency{
		{
			Constructor: func() context.Context {
				return shutdownContext
			},
			Interface: nil,
			Token:     "ShutdownContext",
		},
		{
			Constructor: func() *sync.WaitGroup {
				return &wg
			},
			Interface: nil,
			Token:     "ShutdownWaitGroup",
		},
	})

	go start(container)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	cancel()

	wg.Wait()
}
