package callback_queries

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/constants"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/services"
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
	SenderService       services.ISenderService       `name:"SenderService"`

	ProgramRepository            repositories.IProgramRepository            `name:"ProgramRepository"`
	ExerciseRepository           repositories.IExerciseRepository           `name:"ExerciseRepository"`
	UserExerciseRecordRepository repositories.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
}

type exerciseHandler struct {
	logger                       logger.ILogger
	conversationService          services.IConversationService
	senderService                services.ISenderService
	programRepository            repositories.IProgramRepository
	exerciseRepository           repositories.IExerciseRepository
	userExerciseRecordRepository repositories.IUserExerciseRecordRepository
}

func NewExerciseHandler(deps exerciseHandlerDependencies) *exerciseHandler {
	return &exerciseHandler{
		logger:                       deps.Logger,
		conversationService:          deps.ConversationService,
		senderService:                deps.SenderService,
		programRepository:            deps.ProgramRepository,
		exerciseRepository:           deps.ExerciseRepository,
		userExerciseRecordRepository: deps.UserExerciseRecordRepository,
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

func (h *exerciseHandler) getExerciseName(ctx context.Context, b *tg_bot.Bot, programId uint) string {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	exerciseName := conversation.WaitAnswer()

	existingExercise := h.exerciseRepository.GetByNameAndProgramId(ctx, exerciseName, programId)

	if existingExercise != nil {
		h.senderService.Send(ctx, b, chatId, messages.ExerciseNameAlreadyExistsMessage(exerciseName))
		return h.getExerciseName(ctx, b, programId)
	}

	return exerciseName
}

func (h *exerciseHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	h.senderService.SendSafe(ctx, b, chatId, messages.EnterExerciseNameMessage())

	exerciseName := h.getExerciseName(ctx, b, program.Id)

	exerciseId := h.exerciseRepository.Create(ctx, models.Exercise{
		Name:      exerciseName,
		ProgramId: program.Id,
	})

	userPrograms := h.programRepository.GetAllByProgramId(ctx, program.Id)

	records := make([]models.UserExerciseRecord, 0, 4*len(userPrograms))

	for _, userProgram := range userPrograms {
		for _, rep := range constants.RepsList {
			records = append(records, models.UserExerciseRecord{
				UserProgramId: userProgram.Id,
				ExerciseId:    exerciseId,
				Weight:        0,
				Reps:          uint(rep),
			})
		}
	}

	h.userExerciseRecordRepository.CreateMany(ctx, records)

	msg := messages.ExerciseSuccessfullyAddedMessage(exerciseName, program.Name)
	kb := inline_keyboards.ExerciseOk(program.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *exerciseHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	kb := inline_keyboards.ExerciseOk(program.Id)

	if len(program.Exercises) == 0 {
		msg := messages.NoExercisesMessage(program.Name)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages.ExercisesMessage(program.Name, program.Exercises)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *exerciseHandler) delete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	exercises := h.exerciseRepository.GetByProgramId(ctx, program.Id, limit, offset)

	if len(exercises) == 0 {
		msg := messages.NoExercisesMessage(program.Name)
		kb := inline_keyboards.ExerciseOk(program.Id)

		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	exercisesCount := h.exerciseRepository.CountByProgramId(ctx, program.Id)

	msg := messages.ExerciseDeleteMessage(program.Name)
	kb := inline_keyboards.ProgramExerciseDeleteList(program.Id, exercises, exercisesCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *exerciseHandler) deleteItem(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)
	exercise := utils_context.GetExerciseFromContext(ctx)

	if exercise.ProgramId != program.Id {
		msg := messages.ExerciseNotFoundMessage(exercise.Id)
		kb := inline_keyboards.ExerciseOk(program.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	h.exerciseRepository.DeleteById(ctx, exercise.Id)

	h.userExerciseRecordRepository.DeleteByExerciseId(ctx, exercise.Id)

	msg := messages.ExerciseSuccessfullyDeletedMessage(exercise.Name, program.Name)

	kb := inline_keyboards.ExerciseOk(program.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}
