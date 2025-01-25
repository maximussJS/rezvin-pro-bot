package callback_queries

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/constants/callback_data"
	"rezvin-pro-bot/src/internal/logger"
	repositories2 "rezvin-pro-bot/src/repositories"
	"rezvin-pro-bot/src/services"
	utils_context2 "rezvin-pro-bot/src/utils/context"
	inline_keyboards2 "rezvin-pro-bot/src/utils/inline_keyboards"
	messages2 "rezvin-pro-bot/src/utils/messages"
	"strings"
)

type IBackHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type backHandlerDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`

	SenderService services.ISenderService `name:"SenderService"`

	UserRepository    repositories2.IUserRepository    `name:"UserRepository"`
	ProgramRepository repositories2.IProgramRepository `name:"ProgramRepository"`
}

type backHandler struct {
	logger            logger.ILogger
	senderService     services.ISenderService
	userRepository    repositories2.IUserRepository
	programRepository repositories2.IProgramRepository
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
	chatId := utils_context2.GetChatIdFromContext(ctx)

	kb := inline_keyboards2.ProgramMenu()

	h.senderService.SendWithKb(ctx, b, chatId, messages2.ProgramMenuMessage(), kb)
}

func (h *backHandler) backToProgramList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	programs := h.programRepository.GetAll(ctx, limit, offset)

	if len(programs) == 0 {
		msg := messages2.NoProgramsMessage()
		kb := inline_keyboards2.ProgramMenuOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	programsCount := h.programRepository.CountAll(ctx)

	kb := inline_keyboards2.ProgramList(programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages2.SelectProgramMessage(), kb)
}

func (h *backHandler) backToPendingUsersList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	users := h.userRepository.GetPendingUsers(ctx, limit, offset)

	if len(users) == 0 {
		msg := messages2.NoPendingUsersMessage()
		kb := inline_keyboards2.MainOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	usersCount := h.userRepository.CountPendingUsers(ctx)

	kb := inline_keyboards2.PendingUsersList(users, usersCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages2.SelectPendingUserMessage(), kb)
}

func (h *backHandler) backToClientList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	clients := h.userRepository.GetClients(ctx, limit, offset)

	if len(clients) == 0 {
		msg := messages2.NoClientsMessage()
		kb := inline_keyboards2.MainOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	clientsCount := h.userRepository.CountClients(ctx)

	kb := inline_keyboards2.ClientList(clients, clientsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages2.SelectClientMessage(), kb)
}
