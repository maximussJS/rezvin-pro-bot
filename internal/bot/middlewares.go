package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	bot_utils "rezvin-pro-bot/utils/bot"
	"runtime"
)

func (bot *bot) answerCallbackQueryMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
		answerResult := bot_utils.MustAnswerCallbackQuery(ctx, b, update)

		if !answerResult {
			bot.logger.Error(fmt.Sprintf("Failed to answer callback query: %s", update.CallbackQuery.ID))
			bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   bot.textService.ErrorMessage(),
			})
			return
		}

		next(ctx, b, update)
	}
}

func (bot *bot) isAdminMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
		chatID := bot_utils.GetChatID(update)
		userId := bot_utils.GetUserID(update)

		user := bot.userRepository.GetById(ctx, userId)

		if user == nil || !user.IsAdmin {
			bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
				ChatID: chatID,
				Text:   bot.textService.AdminOnlyMessage(),
			})
			return
		}

		next(ctx, b, update)
	}
}

func (bot *bot) timeoutMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
		timeoutDuration := bot.config.RequestTimeout()

		childCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		doneCh := make(chan struct{})

		go func() {
			next(childCtx, b, update)
			close(doneCh)
		}()

		select {
		case <-childCtx.Done():
			bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
				ChatID: bot_utils.GetChatID(update),
				Text:   bot.textService.RequestTimeoutMessage(),
			})
			return
		case <-doneCh:
			return
		}
	}
}

func (bot *bot) panicRecoveryMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
		defer func() {
			if err := recover(); err != nil {
				chatID := bot_utils.GetChatID(update)
				stackSize := bot.config.ErrorStackTraceSizeInKb() * 1024
				stack := make([]byte, stackSize)
				length := runtime.Stack(stack, true)
				stack = stack[:length]

				if ctx.Err() != nil {
					return
				}

				bot.logger.Error(fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack))

				b.SendMessage(ctx, &tg_bot.SendMessageParams{
					ChatID:    chatID,
					Text:      bot.textService.ErrorMessage(),
					ParseMode: models.ParseModeMarkdown,
				})
			}
		}()

		next(ctx, b, update)
	}
}
