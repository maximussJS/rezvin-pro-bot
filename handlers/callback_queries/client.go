package callback_queries

import (
	"context"
	"fmt"
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
	validate_data "rezvin-pro-bot/utils/validate"
	"strings"
)

type IClientHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type clientHandlerDependencies struct {
	dig.In

	Logger                       logger.ILogger                             `name:"Logger"`
	ConversationService          services.IConversationService              `name:"ConversationService"`
	SenderService                services.ISenderService                    `name:"SenderService"`
	UserRepository               repositories.IUserRepository               `name:"UserRepository"`
	ProgramRepository            repositories.IProgramRepository            `name:"ProgramRepository"`
	UserProgramRepository        repositories.IUserProgramRepository        `name:"UserProgramRepository"`
	UserExerciseRecordRepository repositories.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
}

type clientHandler struct {
	logger                       logger.ILogger
	conversationService          services.IConversationService
	senderService                services.ISenderService
	userRepository               repositories.IUserRepository
	programRepository            repositories.IProgramRepository
	userProgramRepository        repositories.IUserProgramRepository
	userExerciseRecordRepository repositories.IUserExerciseRecordRepository
}

func NewClientHandler(deps clientHandlerDependencies) *clientHandler {
	return &clientHandler{
		logger:                       deps.Logger,
		conversationService:          deps.ConversationService,
		senderService:                deps.SenderService,
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

	if strings.HasPrefix(callBackQueryData, callback_data.ClientResultModifyList) {
		h.resultModifyList(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientResultModifySelected) {
		h.resultModifySelected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientList) {
		h.list(ctx, b)
		return
	}
}

func (h *clientHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	clients := h.userRepository.GetClients(ctx, limit, offset)

	if len(clients) == 0 {
		msg := messages.NoClientsMessage()
		kb := inline_keyboards.MainOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	clientCount := h.userRepository.CountClients(ctx)

	kb := inline_keyboards.ClientList(clients, clientCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectClientMessage(), kb)
}

func (h *clientHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)

	msg := messages.SelectClientOptionMessage(user.GetPrivateName())

	h.senderService.SendWithKb(ctx, b, chatId, msg, inline_keyboards.ClientSelectedMenu(user.Id))
}

func (h *clientHandler) programSelected(ctx context.Context, b *tg_bot.Bot) {
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

	msg := messages.SelectClientProgramOptionMessage(user.GetPrivateName(), userProgram.Name())
	kb := inline_keyboards.ClientSelectedProgramMenu(user.Id, *userProgram)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) programDelete(ctx context.Context, b *tg_bot.Bot) {
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

	h.userProgramRepository.DeleteById(ctx, userProgram.Id)
	h.userExerciseRecordRepository.DeleteByUserProgramId(ctx, userProgram.Id)

	userMsg := messages.UserProgramUnassignedMessage(userProgram.Name())
	h.senderService.SendSafe(ctx, b, user.Id, userMsg)

	adminMsg := messages.ClientProgramDeletedMessage(user.GetPrivateName(), userProgram.Name())
	adminKb := inline_keyboards.ClientSelectedOk(user.Id)

	h.senderService.SendWithKb(ctx, b, chatId, adminMsg, adminKb)
}

func (h *clientHandler) programList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.userProgramRepository.GetByUserId(ctx, user.Id, limit, offset)

	if len(programs) == 0 {
		msg := messages.NoClientProgramsMessage(user.GetPrivateName())
		kb := inline_keyboards.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	programsCount := h.userProgramRepository.CountAllByUserId(ctx, user.Id)

	msg := messages.SelectClientProgramMessage(user.GetPrivateName())

	kb := inline_keyboards.ClientProgramList(user.Id, programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) programAdd(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	programs := h.programRepository.GetNotAssignedToUser(ctx, user.Id, limit, offset)

	if len(programs) == 0 {
		msg := messages.NoProgramsForClientMessage(user.GetPrivateName())
		kb := inline_keyboards.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	programsCount := h.programRepository.CountNotAssignedToUser(ctx, user.Id)

	msg := messages.SelectClientProgramMessage(user.GetPrivateName())

	kb := inline_keyboards.ProgramForClientList(user.Id, programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) programAssign(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	program := utils_context.GetProgramFromContext(ctx)

	userProgram := h.userProgramRepository.GetByUserIdAndProgramId(ctx, user.Id, program.Id)

	if userProgram != nil {
		msg := messages.ClientProgramAlreadyAssignedMessage(user.GetPrivateName(), program.Name)
		kb := inline_keyboards.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	userProgramId := h.userProgramRepository.Create(ctx, models.UserProgram{
		UserId:    user.Id,
		ProgramId: program.Id,
	})

	records := make([]models.UserExerciseRecord, 0, 4*len(program.Exercises))

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

	userMsg := messages.UserProgramAssignedMessage(program.Name)
	h.senderService.SendSafe(ctx, b, user.Id, userMsg)

	adminMsg := messages.ClientProgramAssignedMessage(user.GetPrivateName(), program.Name)
	adminKb := inline_keyboards.ClientSelectedOk(user.Id)

	h.senderService.SendWithKb(ctx, b, chatId, adminMsg, adminKb)
}

func (h *clientHandler) resultList(ctx context.Context, b *tg_bot.Bot) {
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

	records := h.userExerciseRecordRepository.GetAllByUserProgramId(ctx, userProgram.Id)

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

func (h *clientHandler) resultModifyList(ctx context.Context, b *tg_bot.Bot) {
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

	records := h.userExerciseRecordRepository.GetByUserProgramId(ctx, userProgram.Id, limit, offset)

	if len(records) == 0 {
		msg := messages.NoRecordsForClientProgramMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards.ClientProgramSelectedOk(user.Id, userProgram.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	recordsCount := h.userExerciseRecordRepository.CountAllByUserProgramId(ctx, userProgram.Id)

	msg := messages.ClientProgramResultsModifyMessage(user.GetPrivateName(), userProgram.Name())

	kb := inline_keyboards.ClientProgramResultsModifyList(user.Id, records, recordsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) getWeight(ctx context.Context, b *tg_bot.Bot) int {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	answer := conversation.WaitAnswer()

	weight, err := validate_data.ValidateWeightAnswer(answer)

	if err != nil {
		h.senderService.Send(ctx, b, chatId, err.Error())
		return h.getWeight(ctx, b)
	}

	return weight
}

func (h *clientHandler) resultModifySelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	record := utils_context.GetUserExerciseRecordFromContext(ctx)

	h.senderService.SendSafe(ctx, b, chatId, messages.EnterClientResultMessage(user.GetPrivateName(), record.Name()))

	weight := h.getWeight(ctx, b)

	h.userExerciseRecordRepository.UpdateById(ctx, record.Id, models.UserExerciseRecord{
		Weight: weight,
	})

	msg := messages.ClientProgramResultModifiedMessage(user.GetPrivateName(), record.Name())
	kb := inline_keyboards.ClientProgramSelectedOk(user.Id, record.UserProgramId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}
