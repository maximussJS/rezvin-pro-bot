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
)

type IBackHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type backHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                  `name:"Logger"`
	TextService           services.ITextService           `name:"TextService"`
	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`

	UserRepository    repositories.IUserRepository    `name:"UserRepository"`
	ProgramRepository repositories.IProgramRepository `name:"ProgramRepository"`
}

type backHandler struct {
	logger                logger.ILogger
	textService           services.ITextService
	inlineKeyboardService services.IInlineKeyboardService
	userRepository        repositories.IUserRepository
	programRepository     repositories.IProgramRepository
}

func NewBackHandler(deps backHandlerDependencies) *backHandler {
	return &backHandler{
		logger:                deps.Logger,
		textService:           deps.TextService,
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
		h.backToStart(ctx, b, update)
	case callback_data.BackToProgramMenu:
		h.backToProgramMenu(ctx, b, update)
	case callback_data.BackToProgramList:
		h.backToProgramList(ctx, b, update)
	case callback_data.BackToPendingUsersList:
		h.backToPendingUsersList(ctx, b, update)
	case callback_data.BackToClientList:
		h.backToClientList(ctx, b, update)
	}
}

func (h *backHandler) backToMain(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	userId := bot_utils.GetUserID(update)
	chatId := bot_utils.GetChatID(update)
	firstName := bot_utils.GetFirstName(update)
	lastName := bot_utils.GetLastName(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		name := fmt.Sprintf("%s %s", firstName, lastName)
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:      chatId,
			ReplyMarkup: h.inlineKeyboardService.UserRegister(),
			Text:        h.textService.UserRegisterMessage(name),
			ParseMode:   tg_models.ParseModeMarkdown,
		})

		return
	}

	if user.IsAdmin {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:      chatId,
			ReplyMarkup: h.inlineKeyboardService.AdminMain(),
			Text:        h.textService.AdminMainMessage(),
			ParseMode:   tg_models.ParseModeMarkdown,
		})
	} else {
		if user.IsApproved {
			bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
				ChatID:      chatId,
				ReplyMarkup: h.inlineKeyboardService.UserMenu(),
				Text:        h.textService.UserMenuMessage(user.GetReadableName()),
				ParseMode:   tg_models.ParseModeMarkdown,
			})

		} else {
			if user.IsDeclined {
				bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
					ChatID:    chatId,
					Text:      h.textService.DeclinedUserExistsMessage(),
					ParseMode: tg_models.ParseModeMarkdown,
				})
			} else {
				bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
					ChatID:    chatId,
					Text:      h.textService.UnapprovedUserExistsMessage(),
					ParseMode: tg_models.ParseModeMarkdown,
				})
			}
		}
	}
}

func (h *backHandler) backToStart(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.PressStartMessage(),
		ParseMode: tg_models.ParseModeMarkdown,
	})
}

func (h *backHandler) backToProgramMenu(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.ProgramMenuMessage(),
		ReplyMarkup: h.inlineKeyboardService.ProgramMenu(),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *backHandler) backToProgramList(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	programs := h.programRepository.GetAll(ctx, 5, 0)

	chatId := bot_utils.GetChatID(update)

	if len(programs) == 0 {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.NoProgramsMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectProgramMessage(),
		ReplyMarkup: h.inlineKeyboardService.ProgramList(programs),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *backHandler) backToPendingUsersList(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
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

func (h *backHandler) backToClientList(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	clients := h.userRepository.GetClients(ctx, 5, 0)

	if len(clients) == 0 {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.NoClientsMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectClientMessage(),
		ReplyMarkup: h.inlineKeyboardService.ClientList(clients),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}
