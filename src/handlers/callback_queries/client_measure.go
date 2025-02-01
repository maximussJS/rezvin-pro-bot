package callback_queries

import (
	"context"
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

type IClientMeasureHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type clientMeasureHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                      `name:"Logger"`
	SenderService         services.ISenderService             `name:"SenderService"`
	MeasureRepository     repositories.IMeasureRepository     `name:"MeasureRepository"`
	UserMeasureRepository repositories.IUserMeasureRepository `name:"UserMeasureRepository"`
}

type clientMeasureHandler struct {
	logger                logger.ILogger
	senderService         services.ISenderService
	measureRepository     repositories.IMeasureRepository
	userMeasureRepository repositories.IUserMeasureRepository
}

func NewClientMeasureHandler(deps clientMeasureHandlerDependencies) *clientMeasureHandler {
	return &clientMeasureHandler{
		logger:                deps.Logger,
		senderService:         deps.SenderService,
		measureRepository:     deps.MeasureRepository,
		userMeasureRepository: deps.UserMeasureRepository,
	}
}

func (h *clientMeasureHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, constants.ClientMeasureList) {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientProgramAdd) {
		h.add(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientProgramAssign) {
		h.assign(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientProgramSelected) {
		h.selected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientProgramDelete) {
		h.delete(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown client program callback query data: %s", callBackQueryData))
}

func (h *clientProgramHandler) selected(ctx context.Context, b *tg_bot.Bot) {
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
	kb := inline_keyboards.ClientProgramMenu(user.Id, *userProgram)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientProgramHandler) delete(ctx context.Context, b *tg_bot.Bot) {
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
	h.userResultRepository.DeleteByUserProgramId(ctx, userProgram.Id)

	userMsg := messages.UserProgramUnassignedMessage(userProgram.Name())
	userKb := inline_keyboards.UserMenuOk()
	h.senderService.SendWithKb(ctx, b, user.Id, userMsg, userKb)

	adminMsg := messages.ClientProgramDeletedMessage(user.GetPrivateName(), userProgram.Name())
	adminKb := inline_keyboards.ClientSelectedOk(user.Id)

	h.senderService.SendWithKb(ctx, b, chatId, adminMsg, adminKb)
}

func (h *clientMeasureHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	measures := h.userMeasureRepository.GetByUserId(ctx, user.Id, limit, offset)

	if len(measures) == 0 {
		msg := messages.NoClientMeasureMessage(user.GetPrivateName())
		kb := inline_keyboards.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	measuresCount := h.userMeasureRepository.CountAllByUserId(ctx, user.Id)

	msg := messages.SelectClientMeasureMessage(user.GetPrivateName())

	kb := inline_keyboards.ClientMeasuresList(user.Id, measures, measuresCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientProgramHandler) add(ctx context.Context, b *tg_bot.Bot) {
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

	kb := inline_keyboards.ClientProgramAssignList(user.Id, programs, programsCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientProgramHandler) assign(ctx context.Context, b *tg_bot.Bot) {
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

	records := make([]models.UserResult, 0, 4*len(program.Exercises))

	for _, exercise := range program.Exercises {
		for _, rep := range constants.RepsList {
			records = append(records, models.UserResult{
				UserProgramId: userProgramId,
				ExerciseId:    exercise.Id,
				Weight:        0,
				Reps:          uint(rep),
			})
		}
	}

	h.userResultRepository.CreateMany(ctx, records)

	userMsg := messages.UserProgramAssignedMessage(program.Name)
	userKb := inline_keyboards.UserMenuOk()
	h.senderService.SendWithKb(ctx, b, user.Id, userMsg, userKb)

	adminMsg := messages.ClientProgramAssignedMessage(user.GetPrivateName(), program.Name)
	adminKb := inline_keyboards.ClientSelectedOk(user.Id)

	h.senderService.SendWithKb(ctx, b, chatId, adminMsg, adminKb)
}
