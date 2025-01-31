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
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
	"strings"
)

type IProgramHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type programHandlerDependencies struct {
	dig.In

	Logger              logger.ILogger                 `name:"Logger"`
	ConversationService services2.IConversationService `name:"ConversationService"`
	SenderService       services2.ISenderService       `name:"SenderService"`

	ProgramRepository repositories.IProgramRepository `name:"ProgramRepository"`
}

type programHandler struct {
	logger              logger.ILogger
	conversationService services2.IConversationService
	senderService       services2.ISenderService
	programRepository   repositories.IProgramRepository
}

func NewProgramHandler(deps programHandlerDependencies) *programHandler {
	return &programHandler{
		logger:              deps.Logger,
		senderService:       deps.SenderService,
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

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramMenu) {
		h.menu(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ProgramAdd) {
		h.add(ctx, b)
		return
	}
}

func (h *programHandler) menu(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)

	msg := messages.ProgramMenuMessage()

	h.senderService.SendWithKb(ctx, b, chatId, msg, inline_keyboards.ProgramMenu())
}

func (h *programHandler) getProgramName(ctx context.Context, b *tg_bot.Bot) string {
	chatId := utils_context2.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	programName := conversation.WaitAnswer()

	if strings.TrimSpace(programName) == "" {
		h.senderService.Send(ctx, b, chatId, messages.EmptyMessage())
		return h.getProgramName(ctx, b)
	}

	existingProgram := h.programRepository.GetByName(ctx, programName)

	if existingProgram != nil {
		h.senderService.Send(ctx, b, chatId, messages.ProgramNameAlreadyExistsMessage(programName))
		return h.getProgramName(ctx, b)
	}

	return programName
}

func (h *programHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)

	programMsgId := h.senderService.SendSafe(ctx, b, chatId, messages.EnterProgramNameMessage())

	programName := h.getProgramName(ctx, b)

	programId := h.programRepository.Create(ctx, models.Program{
		Name: programName,
	})

	msg := messages.ProgramSuccessfullyAddedMessage(programName)

	kb := inline_keyboards.ProgramOk(programId)

	h.senderService.Delete(ctx, b, chatId, programMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *programHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	programs := h.programRepository.GetAll(ctx, limit, offset)

	if len(programs) == 0 {
		msg := messages.NoProgramsMessage()
		kb := inline_keyboards.ProgramMenuOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	programsCount := h.programRepository.CountAll(ctx)

	kb := inline_keyboards.ProgramList(programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectProgramMessage(), kb)
}

func (h *programHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	program := utils_context2.GetProgramFromContext(ctx)

	msg := messages.SelectProgramOptionMessage(program.Name)
	kb := inline_keyboards.ProgramSelectedMenu(program.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *programHandler) rename(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	program := utils_context2.GetProgramFromContext(ctx)

	programMsgId := h.senderService.SendSafe(ctx, b, chatId, messages.EnterProgramNameMessage())

	programName := h.getProgramName(ctx, b)

	h.programRepository.UpdateById(ctx, program.Id, models.Program{
		Name: programName,
	})

	msg := messages.ProgramSuccessfullyRenamedMessage(program.Name, programName)

	kb := inline_keyboards.ProgramOk(program.Id)

	h.senderService.Delete(ctx, b, chatId, programMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *programHandler) delete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	program := utils_context2.GetProgramFromContext(ctx)

	h.programRepository.DeleteById(ctx, program.Id)

	msg := messages.ProgramSuccessfullyDeletedMessage(program.Name)
	kb := inline_keyboards.ProgramDeleteOk()

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}
