package callback_queries

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/internal/logger"
	"rezvin-pro-bot/src/repositories"
	"rezvin-pro-bot/src/services"
	"rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
	"strings"
)

type IClientHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type clientHandlerDependencies struct {
	dig.In

	Logger         logger.ILogger               `name:"Logger"`
	SenderService  services.ISenderService      `name:"SenderService"`
	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type clientHandler struct {
	logger         logger.ILogger
	senderService  services.ISenderService
	userRepository repositories.IUserRepository
}

func NewClientHandler(deps clientHandlerDependencies) *clientHandler {
	return &clientHandler{
		logger:         deps.Logger,
		senderService:  deps.SenderService,
		userRepository: deps.UserRepository,
	}
}

func (h *clientHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, constants.ClientSelected) {
		h.selected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientList) {
		h.list(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown client callback query data: %s", callBackQueryData))
}

func (h *clientHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	clients := h.userRepository.GetClients(ctx, limit, offset)

	if len(clients) == 0 {
		msg := messages.NoClientsMessage()
		kb := inline_keyboards.MainOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	clientCount := h.userRepository.CountClients(ctx)

	kb := inline_keyboards.ClientList(clients, clientCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectClientMessage(), kb)
}

func (h *clientHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)

	msg := messages.SelectClientOptionMessage(user.GetPrivateName())

	h.senderService.SendWithKb(ctx, b, chatId, msg, inline_keyboards.ClientSelectedMenu(user.Id))
}
