package callback_queries

import (
	"context"
	"errors"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/internal/logger"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/repositories"
	"rezvin-pro-bot/src/services"
	"rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
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

	UserProgramRepository repositories.IUserProgramRepository `name:"UserProgramRepository"`
	ExerciseRepository    repositories.IExerciseRepository    `name:"ExerciseRepository"`
	UserResultRepository  repositories.IUserResultRepository  `name:"UserResultRepository"`
}

type exerciseHandler struct {
	logger                logger.ILogger
	conversationService   services.IConversationService
	senderService         services.ISenderService
	exerciseRepository    repositories.IExerciseRepository
	userProgramRepository repositories.IUserProgramRepository
	userResultRepository  repositories.IUserResultRepository
}

func NewExerciseHandler(deps exerciseHandlerDependencies) *exerciseHandler {
	return &exerciseHandler{
		logger:                deps.Logger,
		conversationService:   deps.ConversationService,
		senderService:         deps.SenderService,
		userProgramRepository: deps.UserProgramRepository,
		exerciseRepository:    deps.ExerciseRepository,
		userResultRepository:  deps.UserResultRepository,
	}
}

func (h *exerciseHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callbackDataQuery := update.CallbackQuery.Data

	if strings.HasPrefix(callbackDataQuery, constants.ExerciseAdd) {
		h.add(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.ExerciseList) {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.ExerciseDeleteItem) {
		h.deleteItem(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.ExerciseDelete) {
		h.delete(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown exercise callback query data: %s", callbackDataQuery))
}

func (h *exerciseHandler) getExerciseName(ctx context.Context, b *tg_bot.Bot, programId uint) (string, error) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	exerciseName := conversation.WaitAnswer()

	if ctx.Err() != nil {
		return "", errors.New("context canceled")
	}

	if strings.TrimSpace(exerciseName) == "" {
		h.senderService.Send(ctx, b, chatId, messages.EmptyMessage())
		return h.getExerciseName(ctx, b, programId)
	}

	existingExercise := h.exerciseRepository.GetByNameAndProgramId(ctx, exerciseName, programId)

	if existingExercise != nil {
		h.senderService.Send(ctx, b, chatId, messages.ExerciseNameAlreadyExistsMessage(exerciseName))
		return h.getExerciseName(ctx, b, programId)
	}

	return exerciseName, nil
}

func (h *exerciseHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	exerciseMsgId := h.senderService.SendSafe(ctx, b, chatId, messages.EnterExerciseNameMessage())

	exerciseName, err := h.getExerciseName(ctx, b, program.Id)

	if err != nil {
		h.senderService.Delete(context.Background(), b, chatId, exerciseMsgId)
		return
	}

	exerciseId := h.exerciseRepository.Create(ctx, models.Exercise{
		Name:      exerciseName,
		ProgramId: program.Id,
	})

	userPrograms := h.userProgramRepository.GetAllByProgramId(ctx, program.Id)

	records := make([]models.UserResult, 0, 4*len(userPrograms))

	for _, userProgram := range userPrograms {
		for _, rep := range constants.RepsList {
			records = append(records, models.UserResult{
				UserProgramId: userProgram.Id,
				ExerciseId:    exerciseId,
				Weight:        0,
				Reps:          uint(rep),
			})
		}
	}

	h.userResultRepository.CreateMany(ctx, records)

	msg := messages.ExerciseSuccessfullyAddedMessage(exerciseName, program.Name)
	kb := inline_keyboards.ExerciseOk(program.Id)
	h.senderService.Delete(ctx, b, chatId, exerciseMsgId)
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
	kb := inline_keyboards.ExerciseDeleteList(program.Id, exercises, exercisesCount, limit, offset)

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

	h.userResultRepository.DeleteByExerciseId(ctx, exercise.Id)

	msg := messages.ExerciseSuccessfullyDeletedMessage(exercise.Name, program.Name)

	kb := inline_keyboards.ExerciseOk(program.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}
