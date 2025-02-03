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
	"rezvin-pro-bot/src/utils/validate"
	"strings"
)

type IClientResultHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type clientResultHandlerDependencies struct {
	dig.In

	Logger               logger.ILogger                     `name:"Logger"`
	ConversationService  services.IConversationService      `name:"ConversationService"`
	SenderService        services.ISenderService            `name:"SenderService"`
	ExerciseRepository   repositories.IExerciseRepository   `name:"ExerciseRepository"`
	UserResultRepository repositories.IUserResultRepository `name:"UserResultRepository"`
}

type clientResultHandler struct {
	logger               logger.ILogger
	conversationService  services.IConversationService
	senderService        services.ISenderService
	exerciseRepository   repositories.IExerciseRepository
	userResultRepository repositories.IUserResultRepository
}

func NewClientResultHandler(deps clientResultHandlerDependencies) *clientResultHandler {
	return &clientResultHandler{
		logger:               deps.Logger,
		conversationService:  deps.ConversationService,
		senderService:        deps.SenderService,
		exerciseRepository:   deps.ExerciseRepository,
		userResultRepository: deps.UserResultRepository,
	}
}

func (h *clientResultHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, constants.ClientResultList) {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientResultExercisesList) {
		h.exerciseList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientResultExerciseSelected) {
		h.exerciseSelected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientResultExerciseReps) {
		h.exerciseRepsSelected(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown client result callback query data: %s", callBackQueryData))
}

func (h *clientResultHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("UserProgram %d not assigned for user %d", userProgram.Id, user.Id))
		msg := messages.ClientProgramNotAssignedMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	records := h.userResultRepository.GetAllByUserProgramId(ctx, userProgram.Id)

	if len(records) == 0 {
		msg := messages.NoRecordsForClientProgramMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards.ClientProgramSelectedOk(user.Id, userProgram.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages.ClientProgramResultsMessage(user.GetPrivateName(), userProgram.Name(), records)

	kb := inline_keyboards.ClientProgramSelectedOk(user.Id, userProgram.Id)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientResultHandler) exerciseList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("UserProgram %d not assigned for user %d", userProgram.Id, user.Id))
		msg := messages.ClientProgramNotAssignedMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	exercises := h.exerciseRepository.GetByProgramId(ctx, userProgram.ProgramId, limit, offset)

	if len(exercises) == 0 {
		msg := messages.NoExercisesMessage(userProgram.Name())
		kb := inline_keyboards.ClientProgramSelectedOk(user.Id, userProgram.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	exercisesCount := h.exerciseRepository.CountByProgramId(ctx, userProgram.ProgramId)

	msg := messages.ClientProgramResultsSelectExerciseMessage(user.GetPrivateName(), userProgram.Name())

	kb := inline_keyboards.ClientResultsExerciseList(user.Id, userProgram.Id, exercises, exercisesCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientResultHandler) getWeight(ctx context.Context, b *tg_bot.Bot) (int, error) {
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

func (h *clientResultHandler) exerciseSelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)
	exercise := utils_context.GetExerciseFromContext(ctx)

	records := h.userResultRepository.GetAllByUserProgramIdAndExerciseId(ctx, userProgram.Id, exercise.Id)

	if len(records) == 0 {
		msg := messages.NoRecordsForClientProgramMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards.ClientProgramSelectedOk(user.Id, userProgram.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages.ClientProgramResultExerciseSelectedMessage(user.GetPrivateName(), exercise.Name)
	kb := inline_keyboards.ClientResultExerciseSelectedOk(user.Id, records)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientResultHandler) exerciseRepsSelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	record := utils_context.GetUserResultFromContext(ctx)

	userMsg := messages.EnterClientResultMessage(user.GetPrivateName(), record.Name())

	userMsgId := h.senderService.SendSafe(ctx, b, chatId, userMsg)

	weight, err := h.getWeight(ctx, b)

	if err != nil {
		h.senderService.Delete(context.Background(), b, chatId, userMsgId)
		return
	}

	h.userResultRepository.UpdateById(ctx, record.Id, models.UserResult{
		Weight: weight,
	})

	msg := messages.ClientProgramResultModifiedMessage(user.GetPrivateName(), record.Name(), record.Reps)
	kb := inline_keyboards.ClientProgramSelectedOk(user.Id, record.UserProgramId)
	h.senderService.Delete(ctx, b, chatId, userMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}
