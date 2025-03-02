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
	bot_utils "rezvin-pro-bot/src/utils/bot"
	"rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
	"strings"
)

type IMainHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type mainHandlerDependencies struct {
	dig.In

	Logger        logger.ILogger          `name:"Logger"`
	SenderService services.ISenderService `name:"SenderService"`

	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type mainHandler struct {
	logger         logger.ILogger
	senderService  services.ISenderService
	userRepository repositories.IUserRepository
}

func NewMainHandler(deps mainHandlerDependencies) *mainHandler {
	return &mainHandler{
		logger:         deps.Logger,
		senderService:  deps.SenderService,
		userRepository: deps.UserRepository,
	}
}

func (h *mainHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, constants.MainBackToMain) {
		h.backToMain(ctx, b, update)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.MainBackToStart) {
		h.backToStart(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown main callback query data: %s", callBackQueryData))
}

func (h *mainHandler) backToMain(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	userId := bot_utils.GetUserID(update)
	firstName := bot_utils.GetFirstName(update)
	lastName := bot_utils.GetLastName(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		name := fmt.Sprintf("%s %s", firstName, lastName)

		kb := inline_keyboards.UserRegister()
		h.senderService.SendWithKb(ctx, b, chatId, messages.NeedRegister(name), kb)
		return
	}

	if user.IsAdmin {
		kb := inline_keyboards.AdminMain()
		h.senderService.SendWithKb(ctx, b, chatId, messages.AdminMainMessage(), kb)
	} else {
		if user.IsApproved {
			msg := messages.UserMenuMessage(user.GetPublicName())
			h.senderService.SendWithKb(ctx, b, chatId, msg, inline_keyboards.UserMenu())
		} else {
			if user.IsDeclined {
				h.senderService.Send(ctx, b, chatId, messages.UserDeclinedMessage(user.GetPublicName()))
			} else {
				h.senderService.Send(ctx, b, chatId, messages.AlreadyRegistered())
			}
		}
	}
}

func (h *mainHandler) backToStart(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	h.senderService.Send(ctx, b, chatId, messages.PressStartMessage())
}
