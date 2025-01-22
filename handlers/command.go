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
)

type ICommandHandler interface {
	Start(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type commandHandlerDependencies struct {
	dig.In

	TextService           services.ITextService           `name:"TextService"`
	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`

	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type commandHandler struct {
	textService           services.ITextService
	inlineKeyboardService services.IInlineKeyboardService

	userRepository repositories.IUserRepository
}

func NewCommandHandler(deps commandHandlerDependencies) *commandHandler {
	return &commandHandler{
		textService:           deps.TextService,
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

	if user == nil {
		name := fmt.Sprintf("%s %s", firstName, lastName)

		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:      chatId,
			ReplyMarkup: c.inlineKeyboardService.UserRegister(),
			Text:        c.textService.UserRegisterMessage(name),
			ParseMode:   models.ParseModeMarkdown,
		})

		return
	}

	if user.IsAdmin {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:      chatId,
			ReplyMarkup: c.inlineKeyboardService.AdminMain(),
			Text:        c.textService.AdminMainMessage(),
			ParseMode:   models.ParseModeMarkdown,
		})
	} else {
		if user.IsApproved {
			bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
				ChatID:      chatId,
				ReplyMarkup: c.inlineKeyboardService.UserMenu(),
				Text:        c.textService.UserMenuMessage(user.GetReadableName()),
				ParseMode:   models.ParseModeMarkdown,
			})

		} else {
			if user.IsDeclined {
				bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
					ChatID:    chatId,
					Text:      c.textService.DeclinedUserExistsMessage(),
					ParseMode: models.ParseModeMarkdown,
				})
			} else {
				bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
					ChatID:    chatId,
					Text:      c.textService.UnapprovedUserExistsMessage(),
					ParseMode: models.ParseModeMarkdown,
				})
			}
		}
	}
}
