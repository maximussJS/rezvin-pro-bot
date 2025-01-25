package callback_queries

import (
	"context"
	"fmt"
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
	inline_keyboards2 "rezvin-pro-bot/src/utils/inline_keyboards"
	messages2 "rezvin-pro-bot/src/utils/messages"
	"rezvin-pro-bot/src/utils/validate"
	"strings"
)

type IClientHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type clientHandlerDependencies struct {
	dig.In

	Logger                       logger.ILogger                              `name:"Logger"`
	ConversationService          services2.IConversationService              `name:"ConversationService"`
	SenderService                services2.ISenderService                    `name:"SenderService"`
	UserRepository               repositories2.IUserRepository               `name:"UserRepository"`
	ProgramRepository            repositories2.IProgramRepository            `name:"ProgramRepository"`
	UserProgramRepository        repositories2.IUserProgramRepository        `name:"UserProgramRepository"`
	UserExerciseRecordRepository repositories2.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
}

type clientHandler struct {
	logger                       logger.ILogger
	conversationService          services2.IConversationService
	senderService                services2.ISenderService
	userRepository               repositories2.IUserRepository
	programRepository            repositories2.IProgramRepository
	userProgramRepository        repositories2.IUserProgramRepository
	userExerciseRecordRepository repositories2.IUserExerciseRecordRepository
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
	chatId := utils_context2.GetChatIdFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	clients := h.userRepository.GetClients(ctx, limit, offset)

	if len(clients) == 0 {
		msg := messages2.NoClientsMessage()
		kb := inline_keyboards2.MainOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	clientCount := h.userRepository.CountClients(ctx)

	kb := inline_keyboards2.ClientList(clients, clientCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages2.SelectClientMessage(), kb)
}

func (h *clientHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)

	msg := messages2.SelectClientOptionMessage(user.GetPrivateName())

	h.senderService.SendWithKb(ctx, b, chatId, msg, inline_keyboards2.ClientSelectedMenu(user.Id))
}

func (h *clientHandler) programSelected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)
	userProgram := utils_context2.GetUserProgramFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("UserProgram %d not assigned for user %d", userProgram.Id, user.Id))
		msg := messages2.ClientProgramNotAssignedMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards2.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages2.SelectClientProgramOptionMessage(user.GetPrivateName(), userProgram.Name())
	kb := inline_keyboards2.ClientSelectedProgramMenu(user.Id, *userProgram)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) programDelete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)
	userProgram := utils_context2.GetUserProgramFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("UserProgram %d not assigned for user %d", userProgram.Id, user.Id))
		msg := messages2.ClientProgramNotAssignedMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards2.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	h.userProgramRepository.DeleteById(ctx, userProgram.Id)
	h.userExerciseRecordRepository.DeleteByUserProgramId(ctx, userProgram.Id)

	userMsg := messages2.UserProgramUnassignedMessage(userProgram.Name())
	userKb := inline_keyboards2.UserMenuOk()
	h.senderService.SendWithKb(ctx, b, user.Id, userMsg, userKb)

	adminMsg := messages2.ClientProgramDeletedMessage(user.GetPrivateName(), userProgram.Name())
	adminKb := inline_keyboards2.ClientSelectedOk(user.Id)

	h.senderService.SendWithKb(ctx, b, chatId, adminMsg, adminKb)
}

func (h *clientHandler) programList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	programs := h.userProgramRepository.GetByUserId(ctx, user.Id, limit, offset)

	if len(programs) == 0 {
		msg := messages2.NoClientProgramsMessage(user.GetPrivateName())
		kb := inline_keyboards2.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	programsCount := h.userProgramRepository.CountAllByUserId(ctx, user.Id)

	msg := messages2.SelectClientProgramMessage(user.GetPrivateName())

	kb := inline_keyboards2.ClientProgramList(user.Id, programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) programAdd(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	programs := h.programRepository.GetNotAssignedToUser(ctx, user.Id, limit, offset)

	if len(programs) == 0 {
		msg := messages2.NoProgramsForClientMessage(user.GetPrivateName())
		kb := inline_keyboards2.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	programsCount := h.programRepository.CountNotAssignedToUser(ctx, user.Id)

	msg := messages2.SelectClientProgramMessage(user.GetPrivateName())

	kb := inline_keyboards2.ProgramForClientList(user.Id, programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) programAssign(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)
	program := utils_context2.GetProgramFromContext(ctx)

	userProgram := h.userProgramRepository.GetByUserIdAndProgramId(ctx, user.Id, program.Id)

	if userProgram != nil {
		msg := messages2.ClientProgramAlreadyAssignedMessage(user.GetPrivateName(), program.Name)
		kb := inline_keyboards2.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	userProgramId := h.userProgramRepository.Create(ctx, models2.UserProgram{
		UserId:    user.Id,
		ProgramId: program.Id,
	})

	records := make([]models2.UserExerciseRecord, 0, 4*len(program.Exercises))

	for _, exercise := range program.Exercises {
		for _, rep := range constants.RepsList {
			records = append(records, models2.UserExerciseRecord{
				UserProgramId: userProgramId,
				ExerciseId:    exercise.Id,
				Weight:        0,
				Reps:          uint(rep),
			})
		}
	}

	h.userExerciseRecordRepository.CreateMany(ctx, records)

	userMsg := messages2.UserProgramAssignedMessage(program.Name)
	userKb := inline_keyboards2.UserMenuOk()
	h.senderService.SendWithKb(ctx, b, user.Id, userMsg, userKb)

	adminMsg := messages2.ClientProgramAssignedMessage(user.GetPrivateName(), program.Name)
	adminKb := inline_keyboards2.ClientSelectedOk(user.Id)

	h.senderService.SendWithKb(ctx, b, chatId, adminMsg, adminKb)
}

func (h *clientHandler) resultList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)
	userProgram := utils_context2.GetUserProgramFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("UserProgram %d not assigned for user %d", userProgram.Id, user.Id))
		msg := messages2.ClientProgramNotAssignedMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards2.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	records := h.userExerciseRecordRepository.GetAllByUserProgramId(ctx, userProgram.Id)

	if len(records) == 0 {
		msg := messages2.NoRecordsForClientProgramMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards2.ClientProgramSelectedOk(user.Id, userProgram.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages2.ClientProgramResultsMessage(user.GetPrivateName(), userProgram.Name(), records)

	kb := inline_keyboards2.ClientProgramSelectedOk(user.Id, userProgram.Id)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) resultModifyList(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)
	userProgram := utils_context2.GetUserProgramFromContext(ctx)
	limit := utils_context2.GetLimitFromContext(ctx)
	offset := utils_context2.GetOffsetFromContext(ctx)

	if userProgram.UserId != user.Id {
		h.logger.Error(fmt.Sprintf("UserProgram %d not assigned for user %d", userProgram.Id, user.Id))
		msg := messages2.ClientProgramNotAssignedMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards2.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	records := h.userExerciseRecordRepository.GetByUserProgramId(ctx, userProgram.Id, limit, offset)

	if len(records) == 0 {
		msg := messages2.NoRecordsForClientProgramMessage(user.GetPrivateName(), userProgram.Name())
		kb := inline_keyboards2.ClientProgramSelectedOk(user.Id, userProgram.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	recordsCount := h.userExerciseRecordRepository.CountAllByUserProgramId(ctx, userProgram.Id)

	msg := messages2.ClientProgramResultsModifyMessage(user.GetPrivateName(), userProgram.Name())

	kb := inline_keyboards2.ClientProgramResultsModifyList(user.Id, records, recordsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientHandler) getWeight(ctx context.Context, b *tg_bot.Bot) int {
	chatId := utils_context2.GetChatIdFromContext(ctx)

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
	chatId := utils_context2.GetChatIdFromContext(ctx)
	user := utils_context2.GetUserFromContext(ctx)
	record := utils_context2.GetUserExerciseRecordFromContext(ctx)

	userMsg := messages2.EnterClientResultMessage(user.GetPrivateName(), record.Name())

	userMsgId := h.senderService.SendSafe(ctx, b, chatId, userMsg)

	weight := h.getWeight(ctx, b)

	h.userExerciseRecordRepository.UpdateById(ctx, record.Id, models2.UserExerciseRecord{
		Weight: weight,
	})

	msg := messages2.ClientProgramResultModifiedMessage(user.GetPrivateName(), record.Name())
	kb := inline_keyboards2.ClientProgramSelectedOk(user.Id, record.UserProgramId)
	h.senderService.Delete(ctx, b, chatId, userMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}
