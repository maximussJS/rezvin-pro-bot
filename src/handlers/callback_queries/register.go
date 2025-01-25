package callback_queries

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/constants/callback_data"
	"rezvin-pro-bot/src/internal/logger"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/repositories"
	"rezvin-pro-bot/src/services"
	bot_utils "rezvin-pro-bot/src/utils/bot"
	"rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
)

type IRegisterHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type registerHandlerDependencies struct {
	dig.In

	SenderService services.ISenderService `name:"SenderService"`

	Logger         logger.ILogger               `name:"Logger"`
	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type registerHandler struct {
	logger         logger.ILogger
	senderService  services.ISenderService
	userRepository repositories.IUserRepository
}

func NewRegisterHandler(deps registerHandlerDependencies) *registerHandler {
	return &registerHandler{
		logger:         deps.Logger,
		senderService:  deps.SenderService,
		userRepository: deps.UserRepository,
	}
}

func (h *registerHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	switch update.CallbackQuery.Data {
	case callback_data.UserRegister:
		h.registerUser(ctx, b, update)
	}
}

func (h *registerHandler) registerUser(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := utils_context.GetChatIdFromContext(ctx)

	userId := bot_utils.GetUserID(update)
	firstName := bot_utils.GetFirstName(update)
	lastName := bot_utils.GetLastName(update)

	user := h.userRepository.GetById(ctx, userId)

	if user != nil {
		if user.IsApproved {
			h.senderService.Send(ctx, b, chatId, messages.AlreadyApprovedRegister())
		} else {
			h.senderService.Send(ctx, b, chatId, messages.AlreadyRegistered())
		}
		return
	}

	h.userRepository.Create(ctx, models.User{
		Id:         userId,
		ChatId:     chatId,
		Username:   bot_utils.GetUsername(update),
		FirstName:  firstName,
		LastName:   lastName,
		IsAdmin:    false,
		IsApproved: false,
		IsDeclined: false,
	})

	h.senderService.Send(ctx, b, chatId, messages.SuccessRegister())

	admins := h.userRepository.GetAdminUsers(ctx)

	name := fmt.Sprintf("%s %s", firstName, lastName)

	msg := messages.NewRegister(name)
	kb := inline_keyboards.MainOk()
	for _, admin := range admins {
		h.senderService.SendWithKb(ctx, b, admin.ChatId, msg, kb)
	}
}
