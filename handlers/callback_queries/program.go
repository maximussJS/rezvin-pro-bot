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

type programHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                  `name:"Logger"`
	TextService           services.ITextService           `name:"TextService"`
	ConversationService   services.IConversationService   `name:"ConversationService"`
	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`

	ProgramRepository repositories.IProgramRepository `name:"ProgramRepository"`
}

type programHandler struct {
	logger                logger.ILogger
	textService           services.ITextService
	conversationService   services.IConversationService
	inlineKeyboardService services.IInlineKeyboardService
	programRepository     repositories.IProgramRepository
}

func NewProgramHandler(deps programHandlerDependencies) *programHandler {
	return &programHandler{
		logger:                deps.Logger,
		textService:           deps.TextService,
		inlineKeyboardService: deps.InlineKeyboardService,
		conversationService:   deps.ConversationService,
		programRepository:     deps.ProgramRepository,
	}
}

func (h *programHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callbackDataQuery := update.CallbackQuery.Data

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramSelected) {
		h.selected(ctx, b, update)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramRename) {
		h.rename(ctx, b, update)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramDelete) {
		h.delete(ctx, b, update)
		return
	}

	switch update.CallbackQuery.Data {
	case callback_data.ProgramMenu:
		h.menu(ctx, b, update)
	case callback_data.ProgramAdd:
		h.add(ctx, b, update)
	case callback_data.ProgramList:
		h.list(ctx, b, update)
	}
}

func (h *programHandler) menu(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.ProgramMenuMessage(),
		ReplyMarkup: h.inlineKeyboardService.ProgramMenu(),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *programHandler) add(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
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
}

func (h *programHandler) list(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
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

func (h *programHandler) selected(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		h.logger.Error(fmt.Sprintf("Program with id %d not found", programId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectProgramOptionMessage(program.Name),
		ParseMode:   tg_models.ParseModeMarkdown,
		ReplyMarkup: h.inlineKeyboardService.ProgramSelectedMenu(programId),
	})
}

func (h *programHandler) rename(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		h.logger.Error(fmt.Sprintf("Program with id %d not found", programId))
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
}

func (h *programHandler) delete(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		h.logger.Error(fmt.Sprintf("Program with id %d not found", programId))
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
}
