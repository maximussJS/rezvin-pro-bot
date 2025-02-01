package services

import (
	"context"
	"rezvin-pro-bot/src/types"
	"sync"
)

type IConversationService interface {
	Shutdown(ctx context.Context) error
	CreateConversation(chatId int64) *types.Conversation
	IsConversationExists(chatId int64) bool
	GetConversation(chatId int64) *types.Conversation
	DeleteConversation(chatId int64)
}

type conversationService struct {
	state map[int64]*types.Conversation
	mu    sync.RWMutex
}

func NewConversationService() *conversationService {
	return &conversationService{
		state: make(map[int64]*types.Conversation),
		mu:    sync.RWMutex{},
	}
}

func (s *conversationService) Shutdown(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for chatId, conv := range s.state {
		conv.Close()
		delete(s.state, chatId)
	}

	return nil
}

func (s *conversationService) CreateConversation(chatId int64) *types.Conversation {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.state[chatId]; ok {
		s.state[chatId].Close()
		delete(s.state, chatId)
	}

	s.state[chatId] = types.NewConversation(chatId)

	return s.state[chatId]
}

func (s *conversationService) IsConversationExists(userId int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.state[userId]
	return ok
}

func (s *conversationService) GetConversation(userId int64) *types.Conversation {
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
