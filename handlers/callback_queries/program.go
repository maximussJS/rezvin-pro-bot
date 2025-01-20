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

type IProgramHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type programDependencies struct {
	dig.In

	Logger                logger.ILogger                  `name:"Logger"`
	TextService           services.ITextService           `name:"TextService"`
	ConversationService   services.IConversationService   `name:"ConversationService"`
	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`

	UserRepository    repositories.IUserRepository    `name:"UserRepository"`
	ProgramRepository repositories.IProgramRepository `name:"ProgramRepository"`
}

type programHandler struct {
	logger                logger.ILogger
	textService           services.ITextService
	conversationService   services.IConversationService
	inlineKeyboardService services.IInlineKeyboardService
	userRepository        repositories.IUserRepository
	programRepository     repositories.IProgramRepository
}

func NewProgramHandler(deps programDependencies) *programHandler {
	return &programHandler{
		logger:                deps.Logger,
		textService:           deps.TextService,
		userRepository:        deps.UserRepository,
		inlineKeyboardService: deps.InlineKeyboardService,
		conversationService:   deps.ConversationService,
		programRepository:     deps.ProgramRepository,
	}
}

func (h *programHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	answerResult := bot_utils.MustAnswerCallbackQuery(ctx, b, update)

	if !answerResult {
		h.logger.Error(fmt.Sprintf("Failed to answer callback query: %s", update.CallbackQuery.ID))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   h.textService.ErrorMessage(),
		})
		return
	}

	if strings.HasPrefix(update.CallbackQuery.Data, callback_data.ProgramSelected) {
		h.programSelected(ctx, b, update)
		return
	}

	if strings.HasPrefix(update.CallbackQuery.Data, callback_data.ProgramRename) {
		h.programRename(ctx, b, update)
		return
	}

	if strings.HasPrefix(update.CallbackQuery.Data, callback_data.ProgramDelete) {
		h.programDelete(ctx, b, update)
		return
	}

	if update.CallbackQuery.Data == callback_data.ProgramBack {
		h.back(ctx, b, update)
	}
}

func (h *programHandler) programSelected(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectProgramOptionMessage(),
		ParseMode:   tg_models.ParseModeMarkdown,
		ReplyMarkup: h.inlineKeyboardService.ProgramSelectedMenu(programId),
	})
}

func (h *programHandler) programRename(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

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

	h.programRepository.UpdateById(ctx, programId, models.Program{
		Name: programName,
	})

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.ProgramSuccessfullyRenamedMessage(program.Name, programName),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	h.back(ctx, b, update)
}

func (h *programHandler) programDelete(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	h.programRepository.DeleteById(ctx, programId)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.ProgramSuccessfullyDeletedMessage(program.Name),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	h.back(ctx, b, update)
}

func (h *programHandler) back(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.AdminProgramMenuMessage(),
		ReplyMarkup: h.inlineKeyboardService.AdminProgramMenu(),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}
