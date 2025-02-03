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
	validate_data "rezvin-pro-bot/src/utils/validate"
	"strings"
)

type IClientMeasureHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type clientMeasureHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                      `name:"Logger"`
	SenderService         services.ISenderService             `name:"SenderService"`
	ConversationService   services.IConversationService       `name:"ConversationService"`
	MeasureRepository     repositories.IMeasureRepository     `name:"MeasureRepository"`
	UserMeasureRepository repositories.IUserMeasureRepository `name:"UserMeasureRepository"`
}

type clientMeasureHandler struct {
	logger                logger.ILogger
	senderService         services.ISenderService
	conversationService   services.IConversationService
	measureRepository     repositories.IMeasureRepository
	userMeasureRepository repositories.IUserMeasureRepository
}

func NewClientMeasureHandler(deps clientMeasureHandlerDependencies) *clientMeasureHandler {
	return &clientMeasureHandler{
		logger:                deps.Logger,
		senderService:         deps.SenderService,
		conversationService:   deps.ConversationService,
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

	if strings.HasPrefix(callBackQueryData, constants.ClientMeasureAdd) {
		h.add(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientMeasureSelected) {
		h.selected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientMeasureDelete) {
		h.delete(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.ClientMeasureResult) {
		h.result(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown client measure callback query data: %s", callBackQueryData))
}

func (h *clientMeasureHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	msg := messages.SelectClientMeasureOptionMessage(user.GetPrivateName(), measure.Name)
	kb := inline_keyboards.ClientMeasureMenu(user.Id, measure.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientMeasureHandler) delete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	lastUserMeasure := h.userMeasureRepository.GetLastByUserIdAndMeasureId(ctx, user.Id, measure.Id)

	if lastUserMeasure == nil {
		msg := messages.NoClientMeasureResultsMessage(user.GetPrivateName(), measure.Name)
		kb := inline_keyboards.ClientMeasureOk(user.Id, measure.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	h.userMeasureRepository.DeleteById(ctx, lastUserMeasure.Id)

	msg := messages.ClientLastMeasureDeletedMessage(user.GetPrivateName(), measure.Name)

	kb := inline_keyboards.ClientMeasureOk(user.Id, measure.Id)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientMeasureHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	measures := h.measureRepository.GetAll(ctx, limit, offset)

	if len(measures) == 0 {
		msg := messages.MeasuresNotFoundMessage()
		kb := inline_keyboards.ClientSelectedOk(user.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	measuresCount := h.measureRepository.CountAll(ctx)

	msg := messages.SelectClientMeasureMessage(user.GetPrivateName())

	kb := inline_keyboards.ClientMeasuresList(user.Id, measures, measuresCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientMeasureHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	msg := messages.EnterClientMeasureValueMessage(user.GetPrivateName(), measure.Name, measure.Units)

	valueMsgId := h.senderService.SendSafe(ctx, b, chatId, msg)

	value, err := h.getValue(ctx, b)

	if err != nil {
		h.senderService.Delete(context.Background(), b, chatId, valueMsgId)
		return
	}

	h.userMeasureRepository.Create(ctx, models.UserMeasure{
		UserId:    user.Id,
		MeasureId: measure.Id,
		Value:     value,
	})

	msg = messages.ClientMeasureAddedMessage(user.GetPrivateName(), measure.Name, measure.Units, value)

	kb := inline_keyboards.ClientMeasureOk(user.Id, measure.Id)

	h.senderService.Delete(ctx, b, chatId, valueMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientMeasureHandler) result(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetUserFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	userMeasures := h.userMeasureRepository.GetAllByUserIdAndMeasureId(ctx, user.Id, measure.Id)

	if len(userMeasures) == 0 {
		msg := messages.NoClientMeasureResultsMessage(user.GetPrivateName(), measure.Name)
		kb := inline_keyboards.ClientMeasureOk(user.Id, measure.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages.ClientMeasureResultMessage(user.GetPrivateName(), *measure, userMeasures)

	kb := inline_keyboards.ClientMeasureOk(user.Id, measure.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *clientMeasureHandler) getValue(ctx context.Context, b *tg_bot.Bot) (float64, error) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	answer := conversation.WaitAnswer()

	if ctx.Err() != nil {
		return 0, errors.New("context canceled")
	}

	value, err := validate_data.ValidateValueAnswer(answer)

	if err != nil {
		h.senderService.Send(ctx, b, chatId, err.Error())
		return h.getValue(ctx, b)
	}

	return value, nil
}
