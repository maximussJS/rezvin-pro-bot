package main

import (
	"context"
	"os"
	"os/signal"
	"rezvin-pro-bot/src/di"
	"rezvin-pro-bot/src/di/dependency"
	"sync"
	"syscall"
)

func main() {
	shutdownContext, cancel := context.WithCancel(context.Background())

	defer cancel()

	var wg sync.WaitGroup

	container := di.BuildContainer()

	container = di.AppendDependenciesToContainer(container, []dependency.Dependency{
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

	go StartApplication(container)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	cancel()

	wg.Wait()
}
