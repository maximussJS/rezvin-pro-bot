package callback_queries

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/services"
	bot_utils "rezvin-pro-bot/utils/bot"
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/inline_keyboards"
	"rezvin-pro-bot/utils/messages"
	"strings"
)

type IPendingUsersHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type pendingUsersHandlerDependencies struct {
	dig.In

	Logger              logger.ILogger                `name:"Logger"`
	ConversationService services.IConversationService `name:"ConversationService"`
	UserRepository      repositories.IUserRepository  `name:"UserRepository"`
}

type pendingUsersHandler struct {
	logger              logger.ILogger
	conversationService services.IConversationService
	userRepository      repositories.IUserRepository
}

func NewPendingUsersHandler(deps pendingUsersHandlerDependencies) *pendingUsersHandler {
	return &pendingUsersHandler{
		logger:              deps.Logger,
		conversationService: deps.ConversationService,
		userRepository:      deps.UserRepository,
	}
}

func (h *pendingUsersHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callbackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callbackQueryData, callback_data.PendingUsersList) {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callbackQueryData, callback_data.PendingUsersSelected) {
		h.selected(ctx, b)
		return
	}

	if strings.HasPrefix(callbackQueryData, callback_data.PendingUsersApprove) {
		h.approve(ctx, b)
		return
	}

	if strings.HasPrefix(callbackQueryData, callback_data.PendingUsersDecline) {
		h.decline(ctx, b)
		return
	}
}

func (h *pendingUsersHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	users := h.userRepository.GetPendingUsers(ctx, limit, offset)

	if len(users) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoPendingUsersMessage())
		return
	}

	usersCount := h.userRepository.CountPendingUsers(ctx)

	kb := inline_keyboards.PendingUsersList(users, usersCount, limit, offset)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.SelectPendingUserMessage(), kb)
}

func (h *pendingUsersHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)

	msg := messages.SelectPendingUserOptionMessage(user.GetPrivateName())
	kb := inline_keyboards.PendingUserDecide(*user)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, kb)
}

func (h *pendingUsersHandler) approve(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)

	h.userRepository.UpdateById(ctx, user.Id, models.User{
		IsApproved: true,
		IsDeclined: false,
	})

	bot_utils.SendMessage(ctx, b, chatId, messages.UserApprovedMessage(user.GetPublicName()))
	bot_utils.SendMessage(ctx, b, chatId, messages.UserApprovedForAdminMessage(user.GetPrivateName()))
}

func (h *pendingUsersHandler) decline(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)

	h.userRepository.UpdateById(ctx, user.Id, models.User{
		IsDeclined: true,
		IsApproved: false,
	})

	bot_utils.SendMessage(ctx, b, chatId, messages.UserDeclinedMessage(user.GetPublicName()))
	bot_utils.SendMessage(ctx, b, chatId, messages.UserDeclinedForAdminMessage(user.GetPrivateName()))
}
