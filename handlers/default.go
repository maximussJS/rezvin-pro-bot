package handlers

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/services"
	bot_utils "rezvin-pro-bot/utils/bot"
)

type IDefaultHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type defaultHandlerDependencies struct {
	dig.In

	TextService         services.ITextService         `name:"TextService"`
	ConversationService services.IConversationService `name:"ConversationService"`
}

type defaultHandler struct {
	textService         services.ITextService
	conversationService services.IConversationService
}

func NewDefaultHandler(deps defaultHandlerDependencies) *defaultHandler {
	return &defaultHandler{
		textService:         deps.TextService,
		conversationService: deps.ConversationService,
	}
}

func (h *defaultHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	chatId := bot_utils.GetChatID(update)

	if h.conversationService.IsConversationExists(chatId) {
		conversation := h.conversationService.GetConversation(chatId)

		conversation.Answer(update.Message.Text)
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.DefaultMessage(),
		ParseMode: models.ParseModeMarkdown,
	})
}
