package callback_queries

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/repositories"
	bot_utils "rezvin-pro-bot/utils/bot"
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/inline_keyboards"
	"rezvin-pro-bot/utils/messages"
	"strings"
)

type IMainHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type mainHandlerDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`

	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type mainHandler struct {
	logger         logger.ILogger
	userRepository repositories.IUserRepository
}

func NewMainHandler(deps mainHandlerDependencies) *mainHandler {
	return &mainHandler{
		logger:         deps.Logger,
		userRepository: deps.UserRepository,
	}
}

func (h *mainHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, callback_data.MainBackToMain) {
		h.backToMain(ctx, b, update)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.MainBackToStart) {
		h.backToStart(ctx, b)
		return
	}
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
		bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.NeedRegister(name), kb)
		return
	}

	if user.IsAdmin {
		kb := inline_keyboards.AdminMain()
		bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.AdminMainMessage(), kb)
	} else {
		if user.IsApproved {
			msg := messages.UserMenuMessage(user.GetPublicName())
			bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, inline_keyboards.UserMenu())
		} else {
			if user.IsDeclined {
				bot_utils.SendMessage(ctx, b, chatId, messages.UserDeclinedMessage(user.GetPublicName()))
			} else {
				bot_utils.SendMessage(ctx, b, chatId, messages.AlreadyRegistered())
			}
		}
	}
}

func (h *mainHandler) backToStart(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	bot_utils.SendMessage(ctx, b, chatId, messages.PressStartMessage())
}
