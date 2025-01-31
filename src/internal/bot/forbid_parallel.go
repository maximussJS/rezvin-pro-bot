package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func (bot *bot) forbidParallel(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		chatId := bot_utils.GetChatID(update)
		userId := bot_utils.GetUserID(update)
		msgId := bot_utils.GetMessageID(update)

		key := fmt.Sprintf("%d:%d", chatId, userId)

		isLocked := bot.lockService.TryLock(key)

		if !isLocked {
			if bot.conversationService.IsConversationExists(chatId) {
				next(ctx, b, update)
				return
			}

			bot.logger.Log(fmt.Sprintf("forbidParallel: chatId: %d, userId: %d,  msgID %d is locked", chatId, userId, msgId))
			return
		}

		defer bot.lockService.Unlock(key)

		next(ctx, b, update)
	}
}
