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

type IMeasureHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type measureHandlerDependencies struct {
	dig.In

	Logger              logger.ILogger                `name:"Logger"`
	ConversationService services.IConversationService `name:"ConversationService"`
	SenderService       services.ISenderService       `name:"SenderService"`

	MeasureRepository repositories.IMeasureRepository `name:"MeasureRepository"`
}

type measureHandler struct {
	logger              logger.ILogger
	conversationService services.IConversationService
	senderService       services.ISenderService
	measureRepository   repositories.IMeasureRepository
}

func NewMeasureHandler(deps measureHandlerDependencies) *measureHandler {
	return &measureHandler{
		logger:              deps.Logger,
		senderService:       deps.SenderService,
		conversationService: deps.ConversationService,
		measureRepository:   deps.MeasureRepository,
	}
}

func (h *measureHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callbackDataQuery := update.CallbackQuery.Data

	if strings.HasPrefix(callbackDataQuery, constants.MeasureMenu) {
		h.menu(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.MeasureAdd) {
		h.add(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.MeasureList) {
		h.list(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.MeasureSelected) {
		h.selected(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.MeasureDelete) {
		h.delete(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.MeasureRename) {
		h.rename(ctx, b)
		return
	}

	if strings.HasPrefix(callbackDataQuery, constants.MeasureChangeUnits) {
		h.changeUnits(ctx, b)
		return
	}

	h.logger.Warn(fmt.Sprintf("Unknown measure callback query data: %s", callbackDataQuery))
}

func (h *measureHandler) menu(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	msg := messages.MeasureMenuMessage()

	h.senderService.SendWithKb(ctx, b, chatId, msg, inline_keyboards.MeasureMenu())
}

func (h *measureHandler) add(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	measureNameMsgId := h.senderService.SendSafe(ctx, b, chatId, messages.EnterMeasureNameMessage())

	measureName, err := h.getMeasureName(ctx, b)

	if err != nil {
		h.senderService.Delete(context.Background(), b, chatId, measureNameMsgId)
		return
	}

	h.senderService.Delete(ctx, b, chatId, measureNameMsgId)
	unitsMsgId := h.senderService.SendSafe(ctx, b, chatId, messages.EnterMeasureUnitsMessage(measureName))

	units, err := h.getMeasureUnits(ctx, b)

	if err != nil {
		h.senderService.Delete(context.Background(), b, chatId, unitsMsgId)
		return
	}

	programId := h.measureRepository.Create(ctx, models.Measure{
		Name:  measureName,
		Units: units,
	})

	msg := messages.MeasureSuccessfullyAddedMessage(measureName, units)

	kb := inline_keyboards.MeasureOk(programId)

	h.senderService.Delete(ctx, b, chatId, unitsMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *measureHandler) list(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	limit := utils_context.GetLimitFromContext(ctx)
	offset := utils_context.GetOffsetFromContext(ctx)

	measures := h.measureRepository.GetAll(ctx, limit, offset)

	if len(measures) == 0 {
		msg := messages.MeasuresNotFoundMessage()
		kb := inline_keyboards.MeasureMenuOk()
		h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
		return
	}

	measuresCount := h.measureRepository.CountAll(ctx)

	kb := inline_keyboards.MeasureList(measures, measuresCount, limit, offset)

	h.senderService.SendWithKb(ctx, b, chatId, messages.SelectMeasureMessage(), kb)
}

func (h *measureHandler) selected(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	msg := messages.SelectMeasureOptionMessage(measure.Name, measure.Units)
	kb := inline_keyboards.MeasureSelectedMenu(measure.Id)

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *measureHandler) delete(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	h.measureRepository.DeleteById(ctx, measure.Id)

	msg := messages.MeasureDeletedMessage(measure.Name)
	kb := inline_keyboards.MeasureDeleteOk()

	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *measureHandler) rename(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	measureNameMsgId := h.senderService.SendSafe(ctx, b, chatId, messages.EnterMeasureNameMessage())

	measureName, err := h.getMeasureName(ctx, b)

	if err != nil {
		h.senderService.Delete(context.Background(), b, chatId, measureNameMsgId)
		return
	}

	h.measureRepository.UpdateById(ctx, measure.Id, models.Measure{
		Name: measureName,
	})

	msg := messages.MeasureRenamed(measure.Name, measureName)
	kb := inline_keyboards.MeasureOk(measure.Id)

	h.senderService.Delete(ctx, b, chatId, measureNameMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *measureHandler) changeUnits(ctx context.Context, b *tg_bot.Bot) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	measure := utils_context.GetMeasureFromContext(ctx)

	unitsMsgId := h.senderService.SendSafe(ctx, b, chatId, messages.EnterMeasureUnitsMessage(measure.Name))

	units, err := h.getMeasureUnits(ctx, b)

	if err != nil {
		h.senderService.Delete(context.Background(), b, chatId, unitsMsgId)
		return
	}

	h.measureRepository.UpdateById(ctx, measure.Id, models.Measure{
		Units: units,
	})

	msg := messages.MeasureUnitsChanged(measure.Name, measure.Units, units)
	kb := inline_keyboards.MeasureOk(measure.Id)

	h.senderService.Delete(ctx, b, chatId, unitsMsgId)
	h.senderService.SendWithKb(ctx, b, chatId, msg, kb)
}

func (h *measureHandler) getMeasureName(ctx context.Context, b *tg_bot.Bot) (string, error) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	measureName := conversation.WaitAnswer()

	if ctx.Err() != nil {
		return "", errors.New("context canceled")
	}

	if strings.TrimSpace(measureName) == "" {
		h.senderService.Send(ctx, b, chatId, messages.EmptyMessage())
		return h.getMeasureName(ctx, b)
	}

	existingMeasure := h.measureRepository.GetByName(ctx, measureName)

	if existingMeasure != nil {
		h.senderService.Send(ctx, b, chatId, messages.MeasureNameAlreadyExistsMessage(measureName))
		return h.getMeasureName(ctx, b)
	}

	return measureName, nil
}

func (h *measureHandler) getMeasureUnits(ctx context.Context, b *tg_bot.Bot) (string, error) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	units := conversation.WaitAnswer()

	if ctx.Err() != nil {
		return "", errors.New("context canceled")
	}

	if strings.TrimSpace(units) == "" {
		h.senderService.Send(ctx, b, chatId, messages.EmptyMessage())
		return h.getMeasureName(ctx, b)
	}

	return units, nil
}
