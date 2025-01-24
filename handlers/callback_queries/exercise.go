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

type IExerciseHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type exerciseHandlerDependencies struct {
	dig.In

	Logger              logger.ILogger                `name:"Logger"`
	ConversationService services.IConversationService `name:"ConversationService"`

	ProgramRepository  repositories.IProgramRepository  `name:"ProgramRepository"`
	ExerciseRepository repositories.IExerciseRepository `name:"ExerciseRepository"`
}

type exerciseHandler struct {
	logger              logger.ILogger
	conversationService services.IConversationService
	programRepository   repositories.IProgramRepository
	exerciseRepository  repositories.IExerciseRepository
}

func NewExerciseHandler(deps exerciseHandlerDependencies) *exerciseHandler {
	return &exerciseHandler{
		logger:              deps.Logger,
		conversationService: deps.ConversationService,
		programRepository:   deps.ProgramRepository,
		exerciseRepository:  deps.ExerciseRepository,
	}
}

func (h *exerciseHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callbackDataQuery := update.CallbackQuery.Data

	if strings.HasPrefix(callbackDataQuery, callback_data.ExerciseAdd) {
		h.add(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ExerciseList) {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ExerciseDeleteItem) {
		h.deleteItem(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ExerciseDelete) {
		h.delete(ctx, b)
		return
	}
}

func (h *exerciseHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	bot_utils.SendMessage(ctx, b, chatId, messages.EnterExerciseNameMessage())

	exerciseName := conversation.WaitAnswer()

	existingExercise := h.exerciseRepository.GetByNameAndProgramId(ctx, exerciseName, program.Id)

	if existingExercise != nil {
		bot_utils.SendMessage(ctx, b, chatId, messages.ExerciseNameAlreadyExistsMessage(exerciseName))
		return
	}

	h.exerciseRepository.Create(ctx, models.Exercise{
		Name:      exerciseName,
		ProgramId: program.Id,
	})

	bot_utils.SendMessage(ctx, b, chatId, messages.ExerciseSuccessfullyAddedMessage(exerciseName, program.Name))

	h.backToSelectedProgram(ctx, b, program)
}

func (h *exerciseHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	if len(program.Exercises) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoExercisesMessage(program.Name))
		return
	}

	bot_utils.SendMessage(ctx, b, chatId, messages.ExercisesMessage(program.Name, program.Exercises))

	h.backToSelectedProgram(ctx, b, program)
}

func (h *exerciseHandler) delete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	exercises := h.exerciseRepository.GetByProgramId(ctx, program.Id, limit, offset)

	if len(exercises) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoExercisesMessage(program.Name))
		return
	}

	exercisesCount := h.exerciseRepository.CountByProgramId(ctx, program.Id)

	msg := messages.ExerciseDeleteMessage(program.Name)
	kb := inline_keyboards.ProgramExerciseDeleteList(program.Id, exercises, exercisesCount, limit, offset)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, kb)
}

func (h *exerciseHandler) deleteItem(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)
	exercise := utils_context.GetExerciseFromContext(ctx)

	if exercise.ProgramId != program.Id {
		bot_utils.SendMessage(ctx, b, chatId, messages.ExerciseNotFoundMessage(exercise.Id))
		return
	}

	h.exerciseRepository.DeleteById(ctx, exercise.Id)

	bot_utils.SendMessage(ctx, b, chatId, messages.ExerciseSuccessfullyDeletedMessage(exercise.Name, program.Name))

	h.backToSelectedProgram(ctx, b, program)
}

func (h *exerciseHandler) backToSelectedProgram(ctx context.Context, b *tg_bot.Bot, program *models.Program) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	msg := messages.SelectProgramOptionMessage(program.Name)
	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, inline_keyboards.ProgramSelectedMenu(program.Id))
}
