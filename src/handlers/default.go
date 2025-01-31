package handlers

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	services2 "rezvin-pro-bot/src/services"
	"rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/messages"
)

type IDefaultHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type defaultHandlerDependencies struct {
	dig.In

	SenderService       services2.ISenderService       `name:"SenderService"`
	ConversationService services2.IConversationService `name:"ConversationService"`
}

type defaultHandler struct {
	senderService       services2.ISenderService
	conversationService services2.IConversationService
}

func NewDefaultHandler(deps defaultHandlerDependencies) *defaultHandler {
	return &defaultHandler{
		senderService:       deps.SenderService,
		conversationService: deps.ConversationService,
	}
}

func (h *defaultHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	if h.conversationService.IsConversationExists(chatId) {
		conversation := h.conversationService.GetConversation(chatId)

		conversation.Answer(update.Message.Text)
		return
	}

	h.senderService.Send(ctx, b, chatId, messages.DefaultMessage())
}
