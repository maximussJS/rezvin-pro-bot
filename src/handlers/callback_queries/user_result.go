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
	utils_context "rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
	"rezvin-pro-bot/src/utils/validate"
	"strings"
)

type IUserResultHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type userResultHandlerDependencies struct {
	dig.In

	Logger               logger.ILogger                     `name:"Logger"`
	ConversationService  services.IConversationService      `name:"ConversationService"`
	SenderService        services.ISenderService            `name:"SenderService"`
	ExerciseRepository   repositories.IExerciseRepository   `name:"ExerciseRepository"`
	UserResultRepository repositories.IUserResultRepository `name:"UserResultRepository"`
}

type userResultHandler struct {
	logger               logger.ILogger
	conversationService  services.IConversationService
	senderService        services.ISenderService
	exerciseRepository   repositories.IExerciseRepository
	userResultRepository repositories.IUserResultRepository
}

func NewUserResultHandler(deps userResultHandlerDependencies) *userResultHandler {
	return &userResultHandler{
		logger:               deps.Logger,
		senderService:        deps.SenderService,
		conversationService:  deps.ConversationService,
		exerciseRepository:   deps.ExerciseRepository,
		userResultRepository: deps.UserResultRepository,
	}
}

func (h *userResultHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, constants.UserResultList) {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.UserResultExerciseList) {
		h.exerciseList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.UserResultExerciseSelected) {
		h.exerciseSelected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.UserResultExerciseReps) {
		h.exerciseRepsSelected(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown user result callback query: %s", callBackQueryData))
}

func (h *userResultHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("Program %d is not assigned for user %d", userProgram.Id, user.Id))
		msg := messages.UserProgramNotAssignedMessage(userProgram.Name())
		kb := inline_keyboards.UserProgramListOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	records := h.userResultRepository.GetAllByUserProgramId(ctx, userProgram.Id)

	kb := inline_keyboards.UserProgramMenuOk(userProgram.Id)

	if len(records) == 0 {
		msg := messages.NoRecordsForUserProgramMessage(userProgram.Name())
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages.UserProgramResultsMessage(userProgram.Name(), records)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userResultHandler) exerciseList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("Program %d is not assigned for user %d", userProgram.Id, user.Id))
		msg := messages.UserProgramNotAssignedMessage(userProgram.Name())
		kb := inline_keyboards.UserProgramListOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	exercises := h.exerciseRepository.GetByProgramId(ctx, userProgram.ProgramId, limit, offset)

	if len(exercises) == 0 {
		msg := messages.NoExercisesMessage(userProgram.Name())
		kb := inline_keyboards.UserProgramMenuOk(userProgram.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	exercisesCount := h.exerciseRepository.CountByProgramId(ctx, userProgram.ProgramId)

	msg := messages.UserProgramResultsSelectExerciseMessage(userProgram.Name())

	kb := inline_keyboards.UserResultExerciseList(userProgram.Id, exercises, exercisesCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userResultHandler) getWeight(ctx context.Context, b *tg_bot.Bot) (int, error) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	answer := conversation.WaitAnswer()

	if ctx.Err() != nil {
		return 0, errors.New("context canceled")
	}

	weight, err := validate_data.ValidateWeightAnswer(answer)

	if err != nil {
		h.senderService.Send(ctx, b, chatId, err.Error())
		return h.getWeight(ctx, b)
	}

	return weight, nil
}

func (h *userResultHandler) exerciseSelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)
	exercise := utils_context.GetExerciseFromContext(ctx)

	records := h.userResultRepository.GetAllByUserProgramIdAndExerciseId(ctx, userProgram.Id, exercise.Id)

	if len(records) == 0 {
		msg := messages.NoRecordsForUserProgramMessage(userProgram.Name())
		kb := inline_keyboards.UserProgramMenuOk(userProgram.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages.UserProgramResultExerciseSelectedMessage(exercise.Name)
	kb := inline_keyboards.UserResultExerciseSelectedOk(records)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userResultHandler) exerciseRepsSelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	record := utils_context.GetUserResultFromContext(ctx)

	userMsg := messages.EnterUserResultMessage(record.Name())

	userMsgId := h.senderService.SendSafe(ctx, b, chatId, userMsg)

	weight, err := h.getWeight(ctx, b)

	if err != nil {
		h.senderService.Delete(context.Background(), b, chatId, userMsgId)
		return
	}

	h.userResultRepository.UpdateById(ctx, record.Id, models.UserResult{
		Weight: weight,
	})

	msg := messages.UserProgramResultModifiedMessage(record.Name(), record.Reps)
	kb := inline_keyboards.UserProgramMenuOk(record.UserProgramId)
	h.senderService.Delete(ctx, b, chatId, userMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}
