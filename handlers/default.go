package handlers

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/services"
	bot_utils "rezvin-pro-bot/utils/bot"
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/messages"
)

type IDefaultHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type defaultHandlerDependencies struct {
	dig.In

	ConversationService services.IConversationService `name:"ConversationService"`
}

type defaultHandler struct {
	conversationService services.IConversationService
}

func NewDefaultHandler(deps defaultHandlerDependencies) *defaultHandler {
	return &defaultHandler{
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

	bot_utils.SendMessage(ctx, b, chatId, messages.DefaultMessage())
}
