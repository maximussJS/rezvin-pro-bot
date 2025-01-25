package handlers

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/services"
	bot_utils "rezvin-pro-bot/utils/bot"
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/inline_keyboards"
	"rezvin-pro-bot/utils/messages"
)

type ICommandHandler interface {
	Start(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type commandHandlerDependencies struct {
	dig.In

	SenderService services.ISenderService `name:"SenderService"`

	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type commandHandler struct {
	senderService  services.ISenderService
	userRepository repositories.IUserRepository
}

func NewCommandHandler(deps commandHandlerDependencies) *commandHandler {
	return &commandHandler{
		senderService:  deps.SenderService,
		userRepository: deps.UserRepository,
	}
}

func (c *commandHandler) Start(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	chatId := utils_context.GetChatIdFromContext(ctx)
	userId := bot_utils.GetUserID(update)
	firstName := bot_utils.GetFirstName(update)
	lastName := bot_utils.GetLastName(update)

	user := c.userRepository.GetById(ctx, userId)

	name := fmt.Sprintf("%s %s", firstName, lastName)

	if user == nil {
		kb := inline_keyboards.UserRegister()
		c.senderService.SendWithKb(ctx, b, chatId, messages.NeedRegister(name), kb)
		return
	}

	if user.IsAdmin {
		kb := inline_keyboards.AdminMain()
		c.senderService.SendWithKb(ctx, b, chatId, messages.AdminMainMessage(), kb)
	} else {
		if user.IsApproved {
			msg := messages.UserMenuMessage(user.GetPublicName())
			c.senderService.SendWithKb(ctx, b, chatId, msg, inline_keyboards.UserMenu())
		} else {
			if user.IsDeclined {
				c.senderService.Send(ctx, b, chatId, messages.UserDeclinedMessage(name))
			} else {
				c.senderService.Send(ctx, b, chatId, messages.AlreadyRegistered())
			}
		}
	}
}
