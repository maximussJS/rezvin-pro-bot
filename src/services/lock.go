package services

import (
	"context"
	"fmt"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/internal/logger"
	"sync"
)

type ILockService interface {
	Lock(key string)
	Unlock(key string)
	TryLock(key string) bool
}

type lockServiceDependencies struct {
	dig.In

	Logger            logger.ILogger  `name:"Logger"`
	ShutdownWaitGroup *sync.WaitGroup `name:"ShutdownWaitGroup"`
	ShutdownContext   context.Context `name:"ShutdownContext"`
}

type lockService struct {
	state             map[string]*sync.RWMutex
	shutdownWaitGroup *sync.WaitGroup
	shutdownContext   context.Context
	logger            logger.ILogger
}

func NewLockService(deps lockServiceDependencies) *lockService {
	s := &lockService{
		state:             make(map[string]*sync.RWMutex),
		shutdownWaitGroup: deps.ShutdownWaitGroup,
		shutdownContext:   deps.ShutdownContext,
		logger:            deps.Logger,
	}

	go s.waitShutdown()

	return s
}

func (ls *lockService) Lock(key string) {
	if _, ok := ls.state[key]; !ok {
		ls.state[key] = &sync.RWMutex{}
	}

	ls.state[key].Lock()
}

func (ls *lockService) Unlock(key string) {
	if _, ok := ls.state[key]; !ok {
		if ls.shutdownContext.Err() != nil {
			return
		}
		panic(fmt.Sprintf("Lock with key %s does not exist", key))
	}

	ls.state[key].Unlock()
}

func (ls *lockService) TryLock(key string) bool {
	if _, ok := ls.state[key]; !ok {
		ls.state[key] = &sync.RWMutex{}
	}

	return ls.state[key].TryLock()
}

func (ls *lockService) waitShutdown() {
	defer ls.shutdownWaitGroup.Done()

	ls.shutdownWaitGroup.Add(1)

	<-ls.shutdownContext.Done()

	for key, _ := range ls.state {
		delete(ls.state, key)
	}

	clear(ls.state)

	ls.logger.Log("Lock service stopped successfully")
}
