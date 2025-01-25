package callback_queries

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/services"
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/inline_keyboards"
	"rezvin-pro-bot/utils/messages"
	"strings"
)

type IBackHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type backHandlerDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`

	SenderService services.ISenderService `name:"SenderService"`

	UserRepository    repositories.IUserRepository    `name:"UserRepository"`
	ProgramRepository repositories.IProgramRepository `name:"ProgramRepository"`
}

type backHandler struct {
	logger            logger.ILogger
	senderService     services.ISenderService
	userRepository    repositories.IUserRepository
	programRepository repositories.IProgramRepository
}

func NewBackHandler(deps backHandlerDependencies) *backHandler {
	return &backHandler{
		logger:            deps.Logger,
		senderService:     deps.SenderService,
		userRepository:    deps.UserRepository,
		programRepository: deps.ProgramRepository,
	}
}

func (h *backHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, callback_data.BackToProgramMenu) {
		h.backToProgramMenu(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.BackToProgramList) {
		h.backToProgramList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.BackToPendingUsersList) {
		h.backToPendingUsersList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.BackToClientList) {
		h.backToClientList(ctx, b)
		return
	}
}

func (h *backHandler) backToProgramMenu(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	kb := inline_keyboards.ProgramMenu()

	h.senderService.SendWithKb(ctx, b, chatId, messages.ProgramMenuMessage(), kb)
}

func (h *backHandler) backToProgramList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.programRepository.GetAll(ctx, limit, offset)

	if len(programs) == 0 {
		h.senderService.Send(ctx, b, chatId, messages.NoProgramsMessage())
		return
	}

	programsCount := h.programRepository.CountAll(ctx)

	kb := inline_keyboards.ProgramList(programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectProgramMessage(), kb)
}

func (h *backHandler) backToPendingUsersList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	users := h.userRepository.GetPendingUsers(ctx, limit, offset)

	if len(users) == 0 {
		h.senderService.Send(ctx, b, chatId, messages.NoPendingUsersMessage())
		return
	}

	usersCount := h.userRepository.CountPendingUsers(ctx)

	kb := inline_keyboards.PendingUsersList(users, usersCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectPendingUserMessage(), kb)
}

func (h *backHandler) backToClientList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	clients := h.userRepository.GetClients(ctx, limit, offset)

	if len(clients) == 0 {
		h.senderService.Send(ctx, b, chatId, messages.NoClientsMessage())
		return
	}

	clientsCount := h.userRepository.CountClients(ctx)

	kb := inline_keyboards.ClientList(clients, clientsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectClientMessage(), kb)
}
