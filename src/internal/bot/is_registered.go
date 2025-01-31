package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	bot_utils "rezvin-pro-bot/src/utils/bot"
	utils_context "rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
)

func (bot *bot) isRegisteredMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		userId := bot_utils.GetUserID(update)

		user := bot.userRepository.GetById(ctx, userId)

		if user == nil {
			chatId := utils_context.GetChatIdFromContext(ctx)
			firstName := bot_utils.GetFirstName(update)
			lastName := bot_utils.GetLastName(update)

			name := fmt.Sprintf("%s %s", firstName, lastName)

			msg := messages.NeedRegister(name)
			kb := inline_keyboards.UserRegister()
			bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
			return
		}

		next(utils_context.GetContextWithCurrentUser(ctx, user), b, update)
	}
}
