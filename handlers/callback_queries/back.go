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
	"rezvin-pro-bot/services"
	bot_utils "rezvin-pro-bot/utils/bot"
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/messages"
)

type IBackHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type backHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                  `name:"Logger"`
	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`

	UserRepository    repositories.IUserRepository    `name:"UserRepository"`
	ProgramRepository repositories.IProgramRepository `name:"ProgramRepository"`
}

type backHandler struct {
	logger                logger.ILogger
	inlineKeyboardService services.IInlineKeyboardService
	userRepository        repositories.IUserRepository
	programRepository     repositories.IProgramRepository
}

func NewBackHandler(deps backHandlerDependencies) *backHandler {
	return &backHandler{
		logger:                deps.Logger,
		inlineKeyboardService: deps.InlineKeyboardService,
		userRepository:        deps.UserRepository,
		programRepository:     deps.ProgramRepository,
	}
}

func (h *backHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	switch update.CallbackQuery.Data {
	case callback_data.BackToMain:
		h.backToMain(ctx, b, update)
	case callback_data.BackToStart:
		h.backToStart(ctx, b)
	case callback_data.BackToProgramMenu:
		h.backToProgramMenu(ctx, b)
	case callback_data.BackToProgramList:
		h.backToProgramList(ctx, b)
	case callback_data.BackToPendingUsersList:
		h.backToPendingUsersList(ctx, b)
	case callback_data.BackToClientList:
		h.backToClientList(ctx, b)
	}
}

func (h *backHandler) backToMain(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	userId := bot_utils.GetUserID(update)
	firstName := bot_utils.GetFirstName(update)
	lastName := bot_utils.GetLastName(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		name := fmt.Sprintf("%s %s", firstName, lastName)

		kb := h.inlineKeyboardService.UserRegister()
		bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.NeedRegister(name), kb)
		return
	}

	if user.IsAdmin {
		kb := h.inlineKeyboardService.AdminMain()
		bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.AdminMainMessage(), kb)
	} else {
		if user.IsApproved {
			msg := messages.UserMenuMessage(user.GetPrivateName())
			bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, h.inlineKeyboardService.UserMenu())
		} else {
			if user.IsDeclined {
				bot_utils.SendMessage(ctx, b, chatId, messages.UserDeclinedMessage(user.GetPublicName()))
			} else {
				bot_utils.SendMessage(ctx, b, chatId, messages.AlreadyRegistered())
			}
		}
	}
}

func (h *backHandler) backToStart(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	bot_utils.SendMessage(ctx, b, chatId, messages.PressStartMessage())
}

func (h *backHandler) backToProgramMenu(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	kb := h.inlineKeyboardService.ProgramMenu()

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.ProgramMenuMessage(), kb)
}

func (h *backHandler) backToProgramList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.programRepository.GetAll(ctx, limit, offset)

	if len(programs) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoProgramsMessage())
		return
	}

	kb := h.inlineKeyboardService.ProgramList(programs)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.SelectProgramMessage(), kb)
}

func (h *backHandler) backToPendingUsersList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	users := h.userRepository.GetPendingUsers(ctx, limit, offset)

	if len(users) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoPendingUsersMessage())
		return
	}

	kb := h.inlineKeyboardService.PendingUsersList(users)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.SelectPendingUserMessage(), kb)
}

func (h *backHandler) backToClientList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	clients := h.userRepository.GetClients(ctx, limit, offset)

	if len(clients) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoClientsMessage())
		return
	}

	kb := h.inlineKeyboardService.ClientList(clients)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.SelectClientMessage(), kb)
}
