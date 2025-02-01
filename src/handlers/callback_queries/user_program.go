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
	utils_context "rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
	"strings"
)

type IUserProgramHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type userProgramHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                      `name:"Logger"`
	SenderService         services.ISenderService             `name:"SenderService"`
	UserProgramRepository repositories.IUserProgramRepository `name:"UserProgramRepository"`
}

type userProgramHandler struct {
	logger                logger.ILogger
	senderService         services.ISenderService
	userProgramRepository repositories.IUserProgramRepository
}

func NewUserProgramHandler(deps userProgramHandlerDependencies) *userProgramHandler {
	return &userProgramHandler{
		logger:                deps.Logger,
		senderService:         deps.SenderService,
		userProgramRepository: deps.UserProgramRepository,
	}
}

func (h *userProgramHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if callBackQueryData == constants.UserProgramList {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.UserProgramSelected) {
		h.selected(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown user program callback query: %s", callBackQueryData))
}

func (h *userProgramHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.userProgramRepository.GetByUserId(ctx, user.Id, limit, offset)

	if len(programs) == 0 {
		msg := messages.NoUserProgramsMessage()
		kb := inline_keyboards.UserMenuOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	programsCount := h.userProgramRepository.CountAllByUserId(ctx, user.Id)

	msg := messages.SelectUserProgramMessage()

	kb := inline_keyboards.UserProgramList(programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userProgramHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("Program %d is not assigned for user %d", userProgram.Id, user.Id))
		msg := messages.UserProgramNotAssignedMessage(userProgram.Name())
		kb := inline_keyboards.UserProgramListOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages.SelectUserProgramOptionMessage(userProgram.Name())
	kb := inline_keyboards.UserProgramMenu(*userProgram)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}
