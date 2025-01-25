package callback_queries

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/constants/callback_data"
	"rezvin-pro-bot/src/internal/logger"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/repositories"
	services2 "rezvin-pro-bot/src/services"
	utils_context2 "rezvin-pro-bot/src/utils/context"
	inline_keyboards2 "rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
	"strings"
)

type IPendingUsersHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type pendingUsersHandlerDependencies struct {
	dig.In

	Logger              logger.ILogger                 `name:"Logger"`
	ConversationService services2.IConversationService `name:"ConversationService"`
	SenderService       services2.ISenderService       `name:"SenderService"`
	UserRepository      repositories.IUserRepository   `name:"UserRepository"`
}

type pendingUsersHandler struct {
	logger              logger.ILogger
	conversationService services2.IConversationService
	senderService       services2.ISenderService
	userRepository      repositories.IUserRepository
}

func NewPendingUsersHandler(deps pendingUsersHandlerDependencies) *pendingUsersHandler {
	return &pendingUsersHandler{
		logger:              deps.Logger,
		senderService:       deps.SenderService,
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
	chatId := utils_context2.GetChatIdFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	users := h.userRepository.GetPendingUsers(ctx, limit, offset)

	if len(users) == 0 {
		msg := messages.NoPendingUsersMessage()
		kb := inline_keyboards2.MainOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	usersCount := h.userRepository.CountPendingUsers(ctx)

	kb := inline_keyboards2.PendingUsersList(users, usersCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectPendingUserMessage(), kb)
}

func (h *pendingUsersHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)

	msg := messages.SelectPendingUserOptionMessage(user.GetPrivateName())
	kb := inline_keyboards2.PendingUserDecide(*user)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *pendingUsersHandler) approve(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)

	h.userRepository.UpdateById(ctx, user.Id, models.User{
		IsApproved: true,
		IsDeclined: false,
	})

	h.senderService.Send(ctx, b, user.ChatId, messages.UserApprovedMessage(user.GetPublicName()))

	adminMsg := messages.UserApprovedForAdminMessage(user.GetPrivateName())
	adminKb := inline_keyboards2.PendingUsersOk()
	h.senderService.SendWithKb(ctx, b, chatId, adminMsg, adminKb)
}

func (h *pendingUsersHandler) decline(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)

	h.userRepository.UpdateById(ctx, user.Id, models.User{
		IsDeclined: true,
		IsApproved: false,
	})

	h.senderService.Send(ctx, b, user.ChatId, messages.UserDeclinedMessage(user.GetPublicName()))

	adminMsg := messages.UserDeclinedForAdminMessage(user.GetPrivateName())
	adminKb := inline_keyboards2.PendingUsersOk()
	h.senderService.SendWithKb(ctx, b, chatId, adminMsg, adminKb)
}
