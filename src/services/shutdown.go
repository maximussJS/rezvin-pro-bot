package services

import (
	"cmp"
	"context"
	"fmt"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/internal/logger"
	"rezvin-pro-bot/src/types"
	"slices"
	"sync"
	"time"
)

type IShutdownService interface {
	Shutdown()
	AddShutdownCallback(callback *types.ShutdownCallback)
}

type shutdownServiceDependencies struct {
	dig.In

	Logger            logger.ILogger  `name:"Logger"`
	ShutdownWaitGroup *sync.WaitGroup `name:"ShutdownWaitGroup"`
	ShutdownContext   context.Context `name:"ShutdownContext"`
}

type shutdownService struct {
	logger            logger.ILogger
	callbacks         []*types.ShutdownCallback
	shutdownWaitGroup *sync.WaitGroup
	shutdownContext   context.Context
}

func NewShutdownService(deps shutdownServiceDependencies) *shutdownService {
	s := &shutdownService{
		shutdownWaitGroup: deps.ShutdownWaitGroup,
		shutdownContext:   deps.ShutdownContext,
		logger:            deps.Logger,
		callbacks:         make([]*types.ShutdownCallback, 0),
	}

	go s.waitForShutdown()

	return s
}

func (s *shutdownService) waitForShutdown() {
	for {
		select {
		case <-s.shutdownContext.Done():
			s.logger.Log("Shutdown signal received")
			s.Shutdown()
			return
		}
	}
}

func (s *shutdownService) AddShutdownCallback(callback *types.ShutdownCallback) {
	s.shutdownWaitGroup.Add(1)
	s.callbacks = append(s.callbacks, callback)
}

func (s *shutdownService) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.logger.Log("Shutting down...")

	sortedCallbacks := make([]*types.ShutdownCallback, len(s.callbacks))

	for i, callback := range s.callbacks {
		sortedCallbacks[i] = callback
	}

	slices.SortFunc(sortedCallbacks, func(a, b *types.ShutdownCallback) int {
		return cmp.Compare(a.Priority, b.Priority)
	})

	for _, callback := range sortedCallbacks {
		err := callback.Callback(ctx)

		if err != nil {
			s.logger.Warn(fmt.Sprintf("Error while executing shutdown for %s callback: %v", callback.Name, err))
		}

		s.logger.Log(fmt.Sprintf("Shutdown callback for %s executed successfully", callback.Name))

		s.shutdownWaitGroup.Done()
	}

	s.logger.Log("Shutdown complete. Bye!")
}
