package callback_queries

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/services"
	bot_utils "rezvin-pro-bot/utils/bot"
	"strings"
)

type IPendingUsersHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type pendingUsersHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                  `name:"Logger"`
	TextService           services.ITextService           `name:"TextService"`
	ConversationService   services.IConversationService   `name:"ConversationService"`
	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`
	UserRepository        repositories.IUserRepository    `name:"UserRepository"`
}

type pendingUsersHandler struct {
	logger                logger.ILogger
	textService           services.ITextService
	conversationService   services.IConversationService
	inlineKeyboardService services.IInlineKeyboardService
	userRepository        repositories.IUserRepository
}

func NewPendingUsersHandler(deps pendingUsersHandlerDependencies) *pendingUsersHandler {
	return &pendingUsersHandler{
		logger:                deps.Logger,
		textService:           deps.TextService,
		conversationService:   deps.ConversationService,
		inlineKeyboardService: deps.InlineKeyboardService,
		userRepository:        deps.UserRepository,
	}
}

func (h *pendingUsersHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	answerResult := bot_utils.MustAnswerCallbackQuery(ctx, b, update)

	if !answerResult {
		h.logger.Error(fmt.Sprintf("Failed to answer callback query: %s", update.CallbackQuery.ID))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   h.textService.ErrorMessage(),
		})
		return
	}

	callbackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callbackQueryData, callback_data.PendingUsersSelected) {
		h.selected(ctx, b, update)
		return
	}

	if strings.HasPrefix(callbackQueryData, callback_data.PendingUsersApprove) {
		h.approve(ctx, b, update)
		return
	}

	if strings.HasPrefix(callbackQueryData, callback_data.PendingUsersDecline) {
		h.decline(ctx, b, update)
		return
	}

	switch callbackQueryData {
	case callback_data.PendingUsersList:
		h.list(ctx, b, update)
	}
}

func (h *pendingUsersHandler) list(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	users := h.userRepository.GetPendingUsers(ctx, 5, 0)

	if len(users) == 0 {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.NoPendingUsersMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectPendingUserMessage(),
		ReplyMarkup: h.inlineKeyboardService.PendingUsersList(users),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *pendingUsersHandler) selected(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	userId := bot_utils.GetSelectedUserId(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		h.logger.Error(fmt.Sprintf("User not found: %d", userId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectPendingUserOptionMessage(fmt.Sprintf("%s %s", user.FirstName, user.LastName)),
		ReplyMarkup: h.inlineKeyboardService.PendingUserDecide(*user),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *pendingUsersHandler) approve(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	userId := bot_utils.GetSelectedUserId(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		h.logger.Error(fmt.Sprintf("User not found: %d", userId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	h.userRepository.UpdateById(ctx, userId, models.User{
		IsApproved: true,
		IsDeclined: false,
	})

	username := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    user.ChatId,
		Text:      h.textService.UserApprovedMessage(username),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.UserApprovedForAdminMessage(username),
		ParseMode: tg_models.ParseModeMarkdown,
	})
}

func (h *pendingUsersHandler) decline(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	userId := bot_utils.GetSelectedUserId(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		h.logger.Error(fmt.Sprintf("User not found: %d", userId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	h.userRepository.UpdateById(ctx, userId, models.User{
		IsDeclined: true,
		IsApproved: false,
	})

	username := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    user.ChatId,
		Text:      h.textService.UserDeclinedMessage(username),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.UserDeclinedForAdminMessage(username),
		ParseMode: tg_models.ParseModeMarkdown,
	})
}
