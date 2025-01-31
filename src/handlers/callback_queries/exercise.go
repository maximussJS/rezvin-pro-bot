package callback_queries

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/constants/callback_data"
	"rezvin-pro-bot/src/internal/logger"
	models2 "rezvin-pro-bot/src/models"
	repositories2 "rezvin-pro-bot/src/repositories"
	services2 "rezvin-pro-bot/src/services"
	utils_context2 "rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
	"strings"
)

type IExerciseHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type exerciseHandlerDependencies struct {
	dig.In

	Logger              logger.ILogger                 `name:"Logger"`
	ConversationService services2.IConversationService `name:"ConversationService"`
	SenderService       services2.ISenderService       `name:"SenderService"`

	ProgramRepository            repositories2.IProgramRepository            `name:"ProgramRepository"`
	UserProgramRepository        repositories2.IUserProgramRepository        `name:"UserProgramRepository"`
	ExerciseRepository           repositories2.IExerciseRepository           `name:"ExerciseRepository"`
	UserExerciseRecordRepository repositories2.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
}

type exerciseHandler struct {
	logger                       logger.ILogger
	conversationService          services2.IConversationService
	senderService                services2.ISenderService
	programRepository            repositories2.IProgramRepository
	exerciseRepository           repositories2.IExerciseRepository
	userProgramRepository        repositories2.IUserProgramRepository
	userExerciseRecordRepository repositories2.IUserExerciseRecordRepository
}

func NewExerciseHandler(deps exerciseHandlerDependencies) *exerciseHandler {
	return &exerciseHandler{
		logger:                       deps.Logger,
		conversationService:          deps.ConversationService,
		senderService:                deps.SenderService,
		programRepository:            deps.ProgramRepository,
		userProgramRepository:        deps.UserProgramRepository,
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
	chatId := utils_context2.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	exerciseName := conversation.WaitAnswer()

	if strings.TrimSpace(exerciseName) == "" {
		h.senderService.Send(ctx, b, chatId, messages.EmptyMessage())
		return h.getExerciseName(ctx, b, programId)
	}

	existingExercise := h.exerciseRepository.GetByNameAndProgramId(ctx, exerciseName, programId)

	if existingExercise != nil {
		h.senderService.Send(ctx, b, chatId, messages.ExerciseNameAlreadyExistsMessage(exerciseName))
		return h.getExerciseName(ctx, b, programId)
	}

	return exerciseName
}

func (h *exerciseHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	program := utils_context2.GetProgramFromContext(ctx)

	exerciseMsgId := h.senderService.SendSafe(ctx, b, chatId, messages.EnterExerciseNameMessage())

	exerciseName := h.getExerciseName(ctx, b, program.Id)

	exerciseId := h.exerciseRepository.Create(ctx, models2.Exercise{
		Name:      exerciseName,
		ProgramId: program.Id,
	})

	userPrograms := h.userProgramRepository.GetAllByProgramId(ctx, program.Id)

	records := make([]models2.UserExerciseRecord, 0, 4*len(userPrograms))

	for _, userProgram := range userPrograms {
		for _, rep := range constants.RepsList {
			records = append(records, models2.UserExerciseRecord{
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
	h.senderService.Delete(ctx, b, chatId, exerciseMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *exerciseHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	program := utils_context2.GetProgramFromContext(ctx)

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
	chatId := utils_context2.GetChatIdFromContext(ctx)
	program := utils_context2.GetProgramFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

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
	chatId := utils_context2.GetChatIdFromContext(ctx)
	program := utils_context2.GetProgramFromContext(ctx)
	exercise := utils_context2.GetExerciseFromContext(ctx)

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
