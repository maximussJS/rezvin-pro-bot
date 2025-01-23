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
	"rezvin-pro-bot/utils/messages"
)

type ICommandHandler interface {
	Start(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type commandHandlerDependencies struct {
	dig.In

	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`

	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type commandHandler struct {
	inlineKeyboardService services.IInlineKeyboardService

	userRepository repositories.IUserRepository
}

func NewCommandHandler(deps commandHandlerDependencies) *commandHandler {
	return &commandHandler{
		userRepository:        deps.UserRepository,
		inlineKeyboardService: deps.InlineKeyboardService,
	}
}

func (c *commandHandler) Start(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	userId := bot_utils.GetUserID(update)
	chatId := bot_utils.GetChatID(update)
	firstName := bot_utils.GetFirstName(update)
	lastName := bot_utils.GetLastName(update)

	user := c.userRepository.GetById(ctx, userId)

	name := fmt.Sprintf("%s %s", firstName, lastName)

	if user == nil {
		kb := c.inlineKeyboardService.UserRegister()
		bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.NeedRegister(name), kb)
		return
	}

	if user.IsAdmin {
		kb := c.inlineKeyboardService.AdminMain()
		bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, messages.AdminMainMessage(), kb)
	} else {
		if user.IsApproved {
			msg := messages.UserMenuMessage(user.GetPrivateName())
			bot_utils.SendMessageWithInlineKeyboard(ctx, b, chatId, msg, c.inlineKeyboardService.UserMenu())
		} else {
			if user.IsDeclined {
				bot_utils.SendMessage(ctx, b, chatId, messages.UserDeclinedMessage(name))
			} else {
				bot_utils.SendMessage(ctx, b, chatId, messages.AlreadyRegistered())
			}
		}
	}
}
