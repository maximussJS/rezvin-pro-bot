package services

import (
	"context"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/internal/logger"
	"sync"
)

type IConversationService interface {
	CreateConversation(chatId int64) *conversation
	IsConversationExists(userId int64) bool
	GetConversation(chatId int64) *conversation
	DeleteConversation(chatId int64)
}

type conversation struct {
	ChatId  int64
	channel chan string
}

func (c *conversation) Close() {
	close(c.channel)
}

func (c *conversation) WaitAnswer() string {
	return <-c.channel
}

func (c *conversation) Answer(text string) {
	c.channel <- text
}

type conversationService struct {
	shutdownWaitGroup *sync.WaitGroup
	shutdownContext   context.Context
	state             map[int64]*conversation
	mu                sync.RWMutex
	logger            logger.ILogger
}

type conversationServiceDependencies struct {
	dig.In

	Logger            logger.ILogger  `name:"Logger"`
	ShutdownWaitGroup *sync.WaitGroup `name:"ShutdownWaitGroup"`
	ShutdownContext   context.Context `name:"ShutdownContext"`
}

func NewConversationService(deps conversationServiceDependencies) *conversationService {
	s := &conversationService{
		shutdownWaitGroup: deps.ShutdownWaitGroup,
		shutdownContext:   deps.ShutdownContext,
		state:             make(map[int64]*conversation),
		mu:                sync.RWMutex{},
		logger:            deps.Logger,
	}

	go s.waitForShutdown()

	return s
}

func (s *conversationService) waitForShutdown() {
	defer s.shutdownWaitGroup.Done()

	s.shutdownWaitGroup.Add(1)

	<-s.shutdownContext.Done()

	s.mu.Lock()
	defer s.mu.Unlock()

	for chatId, conv := range s.state {
		conv.Close()
		delete(s.state, chatId)
	}

	s.logger.Log("Conversation service stopped successfully")
}

func (s *conversationService) CreateConversation(chatId int64) *conversation {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.state[chatId]; ok {
		s.state[chatId].Close()
		delete(s.state, chatId)
	}

	s.state[chatId] = &conversation{
		ChatId:  chatId,
		channel: make(chan string),
	}

	return s.state[chatId]
}

func (s *conversationService) IsConversationExists(userId int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.state[userId]
	return ok
}

func (s *conversationService) GetConversation(userId int64) *conversation {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state[userId]
}

func (s *conversationService) DeleteConversation(userId int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conv, ok := s.state[userId]; ok {
		conv.Close()
		delete(s.state, userId)
	}
}
