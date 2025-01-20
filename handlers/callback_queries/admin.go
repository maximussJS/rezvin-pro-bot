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

type IAdminHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type adminDependencies struct {
	dig.In

	Logger                logger.ILogger                  `name:"Logger"`
	TextService           services.ITextService           `name:"TextService"`
	ConversationService   services.IConversationService   `name:"ConversationService"`
	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`
	UserRepository        repositories.IUserRepository    `name:"UserRepository"`
	ProgramRepository     repositories.IProgramRepository `name:"ProgramRepository"`
}

type adminHandler struct {
	logger                logger.ILogger
	textService           services.ITextService
	conversationService   services.IConversationService
	inlineKeyboardService services.IInlineKeyboardService
	userRepository        repositories.IUserRepository
	programRepository     repositories.IProgramRepository
}

func NewAdminHandler(deps adminDependencies) *adminHandler {
	return &adminHandler{
		logger:                deps.Logger,
		textService:           deps.TextService,
		conversationService:   deps.ConversationService,
		inlineKeyboardService: deps.InlineKeyboardService,
		userRepository:        deps.UserRepository,
		programRepository:     deps.ProgramRepository,
	}
}

func (h *adminHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
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
	case callback_data.AdminProgramMenu:
		h.programMenu(ctx, b, update)
	case callback_data.AdminProgramMenuAdd:
		h.programMenuAdd(ctx, b, update)
	case callback_data.AdminProgramMenuList:
		h.programMenuList(ctx, b, update)
	case callback_data.AdminProgramMenuBack:
		h.backToProgram(ctx, b, update)
	case callback_data.AdminBack:
		h.backToStart(ctx, b, update)
	}
}

func (h *adminHandler) programMenu(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.AdminProgramMenuMessage(),
		ReplyMarkup: h.inlineKeyboardService.AdminProgramMenu(),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *adminHandler) programMenuAdd(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.EnterProgramNameMessage(),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	programName := conversation.WaitAnswer()

	existingProgram := h.programRepository.GetByName(ctx, programName)

	if existingProgram != nil {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ProgramNameAlreadyExistsMessage(programName),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	h.programRepository.Create(ctx, models.Program{
		Name: programName,
	})

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		ParseMode: tg_models.ParseModeMarkdown,
		Text:      h.textService.ProgramSuccessfullyAddedMessage(programName),
	})

	h.backToProgram(ctx, b, update)
}

func (h *adminHandler) programMenuList(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
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
		ReplyMarkup: h.inlineKeyboardService.ProgramMenu(programs),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *adminHandler) backToProgram(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		ReplyMarkup: h.inlineKeyboardService.AdminMenu(),
		Text:        h.textService.AdminStartMessage(),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *adminHandler) backToStart(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.PressStartMessage(),
		ParseMode: tg_models.ParseModeMarkdown,
	})
}
