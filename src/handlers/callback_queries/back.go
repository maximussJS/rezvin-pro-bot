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

type IBackHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type backHandlerDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`

	SenderService services.ISenderService `name:"SenderService"`

	UserRepository    repositories.IUserRepository    `name:"UserRepository"`
	ProgramRepository repositories.IProgramRepository `name:"ProgramRepository"`
	MeasureRepository repositories.IMeasureRepository `name:"MeasureRepository"`
}

type backHandler struct {
	logger            logger.ILogger
	senderService     services.ISenderService
	userRepository    repositories.IUserRepository
	programRepository repositories.IProgramRepository
	measureRepository repositories.IMeasureRepository
}

func NewBackHandler(deps backHandlerDependencies) *backHandler {
	return &backHandler{
		logger:            deps.Logger,
		senderService:     deps.SenderService,
		userRepository:    deps.UserRepository,
		programRepository: deps.ProgramRepository,
		measureRepository: deps.MeasureRepository,
	}
}

func (h *backHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, constants.BackToProgramMenu) {
		h.backToProgramMenu(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.BackToProgramList) {
		h.backToProgramList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.BackToPendingUsersList) {
		h.backToPendingUsersList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.BackToClientList) {
		h.backToClientList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.BackToMeasureList) {
		h.backToMeasureList(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown back callback query data: %s", callBackQueryData))
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
		msg := messages.NoProgramsMessage()
		kb := inline_keyboards.ProgramMenuOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	programsCount := h.programRepository.CountAll(ctx)

	kb := inline_keyboards.ProgramList(programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectProgramMessage(), kb)
}

func (h *backHandler) backToMeasureList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	measures := h.measureRepository.GetAll(ctx, limit, offset)

	if len(measures) == 0 {
		msg := messages.MeasuresNotFoundMessage()
		kb := inline_keyboards.MeasureMenuOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	measuresCount := h.measureRepository.CountAll(ctx)

	kb := inline_keyboards.MeasureList(measures, measuresCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectMeasureMessage(), kb)
}

func (h *backHandler) backToPendingUsersList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	users := h.userRepository.GetPendingUsers(ctx, limit, offset)

	if len(users) == 0 {
		msg := messages.NoPendingUsersMessage()
		kb := inline_keyboards.MainOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
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
		msg := messages.NoClientsMessage()
		kb := inline_keyboards.MainOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	clientsCount := h.userRepository.CountClients(ctx)

	kb := inline_keyboards.ClientList(clients, clientsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectClientMessage(), kb)
}
