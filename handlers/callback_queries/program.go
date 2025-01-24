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

type IProgramHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type programHandlerDependencies struct {
	dig.In

	Logger              logger.ILogger                `name:"Logger"`
	ConversationService services.IConversationService `name:"ConversationService"`

	ProgramRepository repositories.IProgramRepository `name:"ProgramRepository"`
}

type programHandler struct {
	logger              logger.ILogger
	conversationService services.IConversationService
	programRepository   repositories.IProgramRepository
}

func NewProgramHandler(deps programHandlerDependencies) *programHandler {
	return &programHandler{
		logger:              deps.Logger,
		conversationService: deps.ConversationService,
		programRepository:   deps.ProgramRepository,
	}
}

func (h *programHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callbackDataQuery := update.CallbackQuery.Data

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramSelected) {
		h.selected(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramRename) {
		h.rename(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramDelete) {
		h.delete(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramList) {
		h.list(ctx, b)
		return
	}

	switch update.CallbackQuery.Data {
	case callback_data.ProgramMenu:
		h.menu(ctx, b)
	case callback_data.ProgramAdd:
		h.add(ctx, b)
	}
}

func (h *programHandler) menu(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	msg := messages.ProgramMenuMessage()

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, inline_keyboards.ProgramMenu())
}

func (h *programHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	bot_utils.SendMessage(ctx, b, chatId, messages.EnterProgramNameMessage())

	programName := conversation.WaitAnswer()

	existingProgram := h.programRepository.GetByName(ctx, programName)

	if existingProgram != nil {
		bot_utils.SendMessage(ctx, b, chatId, messages.ProgramNameAlreadyExistsMessage(programName))
		return
	}

	h.programRepository.Create(ctx, models.Program{
		Name: programName,
	})

	bot_utils.SendMessage(ctx, b, chatId, messages.ProgramSuccessfullyAddedMessage(programName))
}

func (h *programHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.programRepository.GetAll(ctx, limit, offset)

	if len(programs) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoProgramsMessage())
		return
	}

	programsCount := h.programRepository.CountAll(ctx)

	kb := inline_keyboards.ProgramList(programs, programsCount, limit, offset)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.SelectProgramMessage(), kb)
}

func (h *programHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	msg := messages.SelectProgramOptionMessage(program.Name)
	kb := inline_keyboards.ProgramSelectedMenu(program.Id)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, kb)
}

func (h *programHandler) rename(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	bot_utils.SendMessage(ctx, b, chatId, messages.EnterProgramNameMessage())

	programName := conversation.WaitAnswer()

	existingProgram := h.programRepository.GetByName(ctx, programName)

	if existingProgram != nil {
		bot_utils.SendMessage(ctx, b, chatId, messages.ProgramNameAlreadyExistsMessage(programName))
		return
	}

	h.programRepository.UpdateById(ctx, program.Id, models.Program{
		Name: programName,
	})

	bot_utils.SendMessage(ctx, b, chatId, messages.ProgramSuccessfullyRenamedMessage(program.Name, programName))
}

func (h *programHandler) delete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	h.programRepository.DeleteById(ctx, program.Id)

	bot_utils.SendMessage(ctx, b, chatId, messages.ProgramSuccessfullyDeletedMessage(program.Name))
}
