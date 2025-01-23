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
	bot_utils "rezvin-pro-bot/utils/bot"
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/messages"
	"strings"
)

type IClientHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type clientHandlerDependencies struct {
	dig.In

	Logger                       logger.ILogger                             `name:"Logger"`
	InlineKeyboardService        services.IInlineKeyboardService            `name:"InlineKeyboardService"`
	UserRepository               repositories.IUserRepository               `name:"UserRepository"`
	ProgramRepository            repositories.IProgramRepository            `name:"ProgramRepository"`
	UserProgramRepository        repositories.IUserProgramRepository        `name:"UserProgramRepository"`
	UserExerciseRecordRepository repositories.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
}

type clientHandler struct {
	logger                       logger.ILogger
	inlineKeyboardService        services.IInlineKeyboardService
	userRepository               repositories.IUserRepository
	programRepository            repositories.IProgramRepository
	userProgramRepository        repositories.IUserProgramRepository
	userExerciseRecordRepository repositories.IUserExerciseRecordRepository
}

func NewClientHandler(deps clientHandlerDependencies) *clientHandler {
	return &clientHandler{
		logger:                       deps.Logger,
		inlineKeyboardService:        deps.InlineKeyboardService,
		userRepository:               deps.UserRepository,
		programRepository:            deps.ProgramRepository,
		userProgramRepository:        deps.UserProgramRepository,
		userExerciseRecordRepository: deps.UserExerciseRecordRepository,
	}
}

func (h *clientHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, callback_data.ClientSelected) {
		h.selected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientProgramList) {
		h.programList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientProgramAdd) {
		h.programAdd(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientProgramAssign) {
		h.programAssign(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientProgramSelected) {
		h.programSelected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientProgramDelete) {
		h.programDelete(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientResultList) {
		h.resultList(ctx, b)
		return
	}

	switch callBackQueryData {
	case callback_data.ClientList:
		h.list(ctx, b)
	}
}

func (h *clientHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	clients := h.userRepository.GetClients(ctx, limit, offset)

	if len(clients) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoClientsMessage())
		return
	}

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.SelectClientMessage(), h.inlineKeyboardService.ClientList(clients))
}

func (h *clientHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)

	msg := messages.SelectClientOptionMessage(user.GetPrivateName())

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, h.inlineKeyboardService.ClientSelectedMenu(user.Id))
}

func (h *clientHandler) programSelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)

	msg := messages.SelectClientProgramOptionMessage(user.GetPrivateName(), userProgram.Name())
	kb := h.inlineKeyboardService.ClientSelectedProgramMenu(user.Id, *userProgram)

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) programDelete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)

	h.userProgramRepository.DeleteById(ctx, userProgram.Id)
	h.userExerciseRecordRepository.DeleteByUserProgramId(ctx, userProgram.Id)

	msg := messages.ClientProgramDeletedMessage(user.GetPrivateName(), userProgram.Name())

	bot_utils.SendMessage(ctx, b, chatId, msg)

	h.programList(ctx, b)
}

func (h *clientHandler) programList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.userProgramRepository.GetByUserId(ctx, user.Id, limit, offset)

	if len(programs) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoClientProgramsMessage(user.GetPrivateName()))
		return
	}

	msg := messages.SelectClientProgramMessage(user.GetPrivateName())

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, h.inlineKeyboardService.ClientProgramList(user.Id, programs))
}

func (h *clientHandler) programAdd(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.programRepository.GetNotAssignedToUser(ctx, user.Id, limit, offset)

	if len(programs) == 0 {
		bot_utils.SendMessage(ctx, b, chatId, messages.NoProgramsForClientMessage(user.GetPrivateName()))
		return
	}

	msg := messages.SelectClientProgramMessage(user.GetPrivateName())

	bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, h.inlineKeyboardService.ProgramForClientList(user.Id, programs))
}

func (h *clientHandler) programAssign(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	userProgram := h.userProgramRepository.GetByUserIdAndProgramId(ctx, user.Id, program.Id)

	if userProgram != nil {
		bot_utils.SendMessage(ctx, b, chatId, messages.ClientProgramAlreadyAssignedMessage(user.GetPrivateName(), userProgram.Program.Name))
		return
	}

	userProgramId := h.userProgramRepository.Create(ctx, models.UserProgram{
		UserId:    user.Id,
		ProgramId: program.Id,
	})

	records := make([]models.UserExerciseRecord, 0, len(program.Exercises))

	for _, exercise := range program.Exercises {
		for _, rep := range constants.RepsList {
			records = append(records, models.UserExerciseRecord{
				UserProgramId: userProgramId,
				ExerciseId:    exercise.Id,
				Weight:        0,
				Reps:          uint(rep),
			})
		}
	}

	h.userExerciseRecordRepository.CreateMany(ctx, records)

	bot_utils.SendMessage(ctx, b, chatId, messages.ClientProgramAssignedMessage(user.GetPrivateName(), program.Name))
}

func (h *clientHandler) resultList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	userProgram := utils_context.GetUserProgramFromContext(ctx)

	records := h.userExerciseRecordRepository.GetByUserProgramId(ctx, userProgram.Id)

	if len(records) == 0 {
		msg := messages.NoRecordsForClientProgramMessage(user.GetPrivateName(), userProgram.Name())
		bot_utils.SendMessage(ctx, b, chatId, msg)
		return
	}

	msg := messages.ClientProgramResultsMessage(user.GetPrivateName(), userProgram.Name(), records)

	bot_utils.SendMessage(ctx, b, chatId, msg)
}
