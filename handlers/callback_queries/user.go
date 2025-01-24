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
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/inline_keyboards"
	"rezvin-pro-bot/utils/messages"
	validate_data "rezvin-pro-bot/utils/validate"
	"strings"
)

type IUserHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type userHandlerDependencies struct {
	dig.In

	Logger                       logger.ILogger                             `name:"Logger"`
	ConversationService          services.IConversationService              `name:"ConversationService"`
	UserRepository               repositories.IUserRepository               `name:"UserRepository"`
	UserProgramRepository        repositories.IUserProgramRepository        `name:"UserProgramRepository"`
	UserExerciseRecordRepository repositories.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
}

type userHandler struct {
	logger                       logger.ILogger
	conversationService          services.IConversationService
	userRepository               repositories.IUserRepository
	userProgramRepository        repositories.IUserProgramRepository
	userExerciseRecordRepository repositories.IUserExerciseRecordRepository
}

func NewUserHandler(deps userHandlerDependencies) *userHandler {
	return &userHandler{
		logger:                       deps.Logger,
		conversationService:          deps.ConversationService,
		userRepository:               deps.UserRepository,
		userProgramRepository:        deps.UserProgramRepository,
		userExerciseRecordRepository: deps.UserExerciseRecordRepository,
	}
}

func (h *userHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if callBackQueryData == callback_data.UserProgramList {
		h.programList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.UserProgramSelected) {
		h.programSelected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.UserResultList) {
		h.resultList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.UserResultModifyList) {
		h.resultModifyList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.UserResultModifySelected) {
		h.resultModifySelected(ctx, b)
		return
	}
}

func (h *userHandler) programList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.userProgramRepository.GetByUserId(ctx, user.Id, limit, offset)

	if len(programs) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoUserProgramsMessage())
		return
	}

	programsCount := h.userProgramRepository.CountAllByUserId(ctx, user.Id)

	msg := messages.SelectUserProgramMessage()

	kb := inline_keyboards.UserProgramList(programs, programsCount, limit, offset)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, kb)
}

func (h *userHandler) programSelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("Program %d is not assigned for user %d", userProgram.Id, user.Id))
		bot_utils.SendMessage(ctx, b, chatId, messages.UserProgramNotAssignedMessage(userProgram.Name()))
		return
	}

	msg := messages.SelectUserProgramOptionMessage(userProgram.Name())
	kb := inline_keyboards.UserProgramMenu(*userProgram)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, kb)
}

func (h *userHandler) resultList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("Program %d is not assigned for user %d", userProgram.Id, user.Id))
		bot_utils.SendMessage(ctx, b, chatId, messages.UserProgramNotAssignedMessage(userProgram.Name()))
		return
	}

	records := h.userExerciseRecordRepository.GetAllByUserProgramId(ctx, userProgram.Id)

	if len(records) == 0 {
		msg := messages.NoRecordsForUserProgramMessage(userProgram.Name())
		bot_utils.SendMessage(ctx, b, chatId, msg)
		return
	}

	msg := messages.UserProgramResultsMessage(userProgram.Name(), records)

	bot_utils.SendMessage(ctx, b, chatId, msg)

	h.programSelected(ctx, b)
}

func (h *userHandler) resultModifyList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("Program %d is not assigned for user %d", userProgram.Id, user.Id))
		bot_utils.SendMessage(ctx, b, chatId, messages.UserProgramNotAssignedMessage(userProgram.Name()))
		return
	}

	records := h.userExerciseRecordRepository.GetByUserProgramId(ctx, userProgram.Id, limit, offset)

	if len(records) == 0 {
		msg := messages.NoRecordsForUserProgramMessage(userProgram.Name())
		bot_utils.SendMessage(ctx, b, chatId, msg)
		return
	}

	recordsCount := h.userExerciseRecordRepository.CountAllByUserProgramId(ctx, userProgram.Id)

	msg := messages.UserProgramResultsModifyMessage(userProgram.Name())

	kb := inline_keyboards.UserProgramResultsModifyList(records, recordsCount, limit, offset)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, kb)
}

func (h *userHandler) getWeight(ctx context.Context, b *tg_bot.Bot) int {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	answer := conversation.WaitAnswer()

	weight, err := validate_data.ValidateWeightAnswer(answer)

	if err != nil {
		bot_utils.SendMessage(ctx, b, chatId, err.Error())
		return h.getWeight(ctx, b)
	}

	return weight
}

func (h *userHandler) resultModifySelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	record := utils_context.GetUserExerciseRecordFromContext(ctx)

	bot_utils.SendMessage(ctx, b, chatId, messages.EnterUserResultMessage(record.Name()))

	weight := h.getWeight(ctx, b)

	h.userExerciseRecordRepository.UpdateById(ctx, record.Id, models.UserExerciseRecord{
		Weight: weight,
	})

	msg := messages.UserProgramResultModifiedMessage(record.Name())

	bot_utils.SendMessage(ctx, b, chatId, msg)

	h.programSelected(ctx, b)
}