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
	Shutdown(ctx context.Context) error
}

type lockServiceDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`
}

type lockService struct {
	state  map[string]*sync.RWMutex
	logger logger.ILogger
}

func NewLockService(deps lockServiceDependencies) *lockService {
	return &lockService{
		state:  make(map[string]*sync.RWMutex),
		logger: deps.Logger,
	}
}

func (ls *lockService) Lock(key string) {
	if _, ok := ls.state[key]; !ok {
		ls.state[key] = &sync.RWMutex{}
	}

	ls.state[key].Lock()
}

func (ls *lockService) Unlock(key string) {
	if _, ok := ls.state[key]; !ok {
		ls.logger.Warn(fmt.Sprintf("Lock with key %s does not exist", key))
	}

	ls.state[key].Unlock()
}

func (ls *lockService) TryLock(key string) bool {
	if _, ok := ls.state[key]; !ok {
		ls.state[key] = &sync.RWMutex{}
	}

	return ls.state[key].TryLock()
}

func (ls *lockService) Shutdown(_ context.Context) error {
	for key, _ := range ls.state {
		delete(ls.state, key)
	}

	clear(ls.state)

	return nil
}
