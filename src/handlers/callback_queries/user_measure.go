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

type IUserMeasureHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type userMeasureHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                      `name:"Logger"`
	SenderService         services.ISenderService             `name:"SenderService"`
	ConversationService   services.IConversationService       `name:"ConversationService"`
	MeasureRepository     repositories.IMeasureRepository     `name:"MeasureRepository"`
	UserMeasureRepository repositories.IUserMeasureRepository `name:"UserMeasureRepository"`
}

type userMeasureHandler struct {
	logger                logger.ILogger
	senderService         services.ISenderService
	conversationService   services.IConversationService
	measureRepository     repositories.IMeasureRepository
	userMeasureRepository repositories.IUserMeasureRepository
}

func NewUserMeasureHandler(deps userMeasureHandlerDependencies) *userMeasureHandler {
	return &userMeasureHandler{
		logger:                deps.Logger,
		senderService:         deps.SenderService,
		conversationService:   deps.ConversationService,
		measureRepository:     deps.MeasureRepository,
		userMeasureRepository: deps.UserMeasureRepository,
	}
}

func (h *userMeasureHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, constants.UserMeasureList) {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.UserMeasureAdd) {
		h.add(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.UserMeasureSelected) {
		h.selected(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.UserMeasureDelete) {
		h.delete(ctx, b)
		return
	}

	if strings.HasPrefix(callBackQueryData, constants.UserMeasureResult) {
		h.result(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown user measure callback query data: %s", callBackQueryData))
}

func (h *userMeasureHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	msg := messages.SelectUserMeasureOptionMessage(measure.Name)
	kb := inline_keyboards.UserMeasureMenu(measure.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userMeasureHandler) delete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	lastUserMeasure := h.userMeasureRepository.GetLastByUserIdAndMeasureId(ctx, user.Id, measure.Id)

	if lastUserMeasure == nil {
		msg := messages.NoUserMeasureResultsMessage(measure.Name)
		kb := inline_keyboards.UserMeasureOk(measure.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	h.userMeasureRepository.DeleteById(ctx, lastUserMeasure.Id)

	msg := messages.UserLastMeasureDeletedMessage(measure.Name)

	kb := inline_keyboards.UserMeasureOk(measure.Id)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userMeasureHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	measures := h.measureRepository.GetAll(ctx, limit, offset)

	if len(measures) == 0 {
		msg := messages.MeasuresNotFoundMessage()
		kb := inline_keyboards.UserMenuOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	measuresCount := h.measureRepository.CountAll(ctx)

	msg := messages.SelectUserMeasureMessage()

	kb := inline_keyboards.UserMeasuresList(measures, measuresCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userMeasureHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	msg := messages.EnterUserMeasureValueMessage(measure.Name, measure.Units)

	valueMsgId := h.senderService.SendSafe(ctx, b, chatId, msg)

	value, err := h.getValue(ctx, b)

	if err != nil {
		h.senderService.Delete(ctx, b, chatId, valueMsgId)
		return
	}

	h.userMeasureRepository.Create(ctx, models.UserMeasure{
		UserId:    user.Id,
		MeasureId: measure.Id,
		Value:     value,
	})

	msg = messages.UserMeasureAddedMessage(measure.Name, measure.Units, value)

	kb := inline_keyboards.UserMeasureOk(measure.Id)

	h.senderService.Delete(ctx, b, chatId, valueMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userMeasureHandler) result(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	user := utils_context.GetCurrentUserFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	userMeasures := h.userMeasureRepository.GetAllByUserIdAndMeasureId(ctx, user.Id, measure.Id)

	if len(userMeasures) == 0 {
		msg := messages.NoUserMeasureResultsMessage(measure.Name)
		kb := inline_keyboards.UserMeasureOk(measure.Id)
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	msg := messages.UserMeasureResultMessage(*measure, userMeasures)

	kb := inline_keyboards.UserMeasureOk(measure.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *userMeasureHandler) getValue(ctx context.Context, b *tg_bot.Bot) (float64, error) {
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
