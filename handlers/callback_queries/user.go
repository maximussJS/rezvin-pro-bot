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
)

type IUserHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type userDependencies struct {
	dig.In

	Logger         logger.ILogger               `name:"Logger"`
	TextService    services.ITextService        `name:"TextService"`
	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type userHandler struct {
	logger         logger.ILogger
	textService    services.ITextService
	userRepository repositories.IUserRepository
}

func NewUserHandler(deps userDependencies) *userHandler {
	return &userHandler{
		logger:         deps.Logger,
		textService:    deps.TextService,
		userRepository: deps.UserRepository,
	}
}

func (h *userHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	answerResult := bot_utils.MustAnswerCallbackQuery(ctx, b, update)

	if !answerResult {
		h.logger.Error(fmt.Sprintf("Failed to answer callback query: %s", update.CallbackQuery.ID))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   h.textService.ErrorMessage(),
		})
		return
	}

	switch update.CallbackQuery.Data {
	case callback_data.UserRegister:
		h.registerUser(ctx, b, update)
	case callback_data.UserBack:
		h.back(ctx, b, update)
	}
}

func (h *userHandler) registerUser(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	userId := bot_utils.GetUserID(update)
	chatId := bot_utils.GetChatID(update)

	user := h.userRepository.GetById(ctx, userId)

	if user != nil {
		if user.IsApproved {
			bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
				ChatID:    chatId,
				Text:      h.textService.ApprovedUserExistsMessage(),
				ParseMode: tg_models.ParseModeMarkdown,
			})
		} else {
			bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
				ChatID:    chatId,
				Text:      h.textService.UnapprovedUserExistsMessage(),
				ParseMode: tg_models.ParseModeMarkdown,
			})
		}
		return
	}

	h.userRepository.Create(ctx, models.User{
		Id:         userId,
		ChatId:     bot_utils.GetChatID(update),
		Username:   bot_utils.GetUsername(update),
		FirstName:  bot_utils.GetFirstName(update),
		LastName:   bot_utils.GetLastName(update),
		IsAdmin:    false,
		IsApproved: false,
	})

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.UserRegisterSuccessMessage(),
		ParseMode: tg_models.ParseModeMarkdown,
	})
}

func (h *userHandler) back(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.PressStartMessage(),
		ParseMode: tg_models.ParseModeMarkdown,
	})
}
