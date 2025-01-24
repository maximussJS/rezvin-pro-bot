package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/models"
	bot_utils "rezvin-pro-bot/utils/bot"
	utils_context "rezvin-pro-bot/utils/context"
	"rezvin-pro-bot/utils/messages"
	"runtime"
)

func (bot *bot) answerCallbackQueryMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		answerResult := bot_utils.AnswerCallbackQuery(ctx, b, update)

		if !answerResult {
			chatId := utils_context.GetChatIdFromContext(ctx)

			bot.logger.Error(fmt.Sprintf("Failed to answer callback query: %s", update.CallbackQuery.ID))
			bot_utils.SendMessage(ctx, b, chatId, messages.ErrorMessage())
			return
		}

		next(ctx, b, update)
	}
}

func (bot *bot) isRegisteredMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		userId := bot_utils.GetUserID(update)

		user := bot.userRepository.GetById(ctx, userId)

		if user == nil {
			chatId := utils_context.GetChatIdFromContext(ctx)
			firstName := bot_utils.GetFirstName(update)
			lastName := bot_utils.GetLastName(update)

			name := fmt.Sprintf("%s %s", firstName, lastName)

			bot_utils.SendMessage(ctx, b, chatId, messages.NeedRegister(name))
			return
		}

		if user.IsNotApproved() {
			chatId := utils_context.GetChatIdFromContext(ctx)

			bot_utils.SendMessage(ctx, b, chatId, messages.UserNotApprovedMessage())
			return
		}

		next(utils_context.GetContextWithCurrentUser(ctx, user), b, update)
	}
}

func (bot *bot) isApprovedMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		user := utils_context.GetCurrentUserFromContext(ctx)

		if user.IsNotApproved() {
			chatId := utils_context.GetChatIdFromContext(ctx)

			bot_utils.SendMessage(ctx, b, chatId, messages.UserNotApprovedMessage())
			return
		}

		next(utils_context.GetContextWithCurrentUser(ctx, user), b, update)
	}
}

func (bot *bot) isAdminMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		user := utils_context.GetCurrentUserFromContext(ctx)

		if user.IsNotAdmin() {
			chatId := utils_context.GetChatIdFromContext(ctx)

			bot_utils.SendMessage(ctx, b, chatId, messages.AdminOnlyMessage())
			return
		}

		next(ctx, b, update)
	}
}

func (bot *bot) chatIdMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		chatId := bot_utils.GetChatID(update)
		userId := bot_utils.GetUserID(update)

		user := bot.userRepository.GetById(ctx, userId)

		if user != nil && user.ChatId != chatId {
			user.ChatId = chatId
			bot.userRepository.UpdateById(ctx, userId, models.User{
				ChatId: chatId,
			})
		}

		next(utils_context.GetContextWithChatId(ctx, chatId), b, update)
	}
}

func (bot *bot) parseParamsMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		callbackQueryData := update.CallbackQuery.Data
		chatId := utils_context.GetChatIdFromContext(ctx)

		params, err := bot_utils.ParseParamsFromQueryString(callbackQueryData)

		if err != nil {
			bot.logger.Error(fmt.Sprintf("Failed to parse params: %s", callbackQueryData))
			bot_utils.SendMessage(ctx, b, chatId, messages.ParamsErrorMessage(err))
			return
		}

		if params == nil {
			bot.logger.Error(fmt.Sprintf("Failed to parse params: %s", update.Message.Text))
			bot_utils.SendMessage(ctx, b, chatId, messages.ErrorMessage())
			return
		}

		next(utils_context.GetContextWithParams(ctx, params), b, update)
	}
}

func (bot *bot) validateParamsMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		chatId := utils_context.GetChatIdFromContext(ctx)
		params := utils_context.GetParamsFromContext(ctx)

		newCtx := context.WithoutCancel(ctx)

		if params.UserId != 0 {
			user := bot.userRepository.GetById(ctx, params.UserId)

			if user == nil {
				bot_utils.SendMessage(ctx, b, chatId, messages.UserNotFoundMessage(params.UserId))
				return
			}

			newCtx = utils_context.GetContextWithUser(newCtx, user)
		}

		if params.ProgramId != 0 {
			program := bot.programRepository.GetById(ctx, params.ProgramId)

			if program == nil {
				bot_utils.SendMessage(ctx, b, chatId, messages.ProgramNotFoundMessage(params.ProgramId))
				return
			}

			newCtx = utils_context.GetContextWithProgram(newCtx, program)
		}

		if params.ExerciseId != 0 {
			exercise := bot.exerciseRepository.GetById(ctx, params.ExerciseId)

			if exercise == nil {
				bot_utils.SendMessage(ctx, b, chatId, messages.ExerciseNotFoundMessage(params.ExerciseId))
				return
			}

			newCtx = utils_context.GetContextWithExercise(newCtx, exercise)
		}

		if params.UserProgramId != 0 {
			userProgram := bot.userProgramRepository.GetById(ctx, params.UserProgramId)

			if userProgram == nil {
				bot_utils.SendMessage(ctx, b, chatId, messages.ClientProgramNotFoundMessage(params.UserProgramId))
				return
			}

			newCtx = utils_context.GetContextWithUserProgram(newCtx, userProgram)
		}

		if params.UserExerciseRecordId != 0 {
			record := bot.userExerciseRecordRepository.GetById(ctx, params.UserExerciseRecordId)

			if record == nil {
				bot_utils.SendMessage(ctx, b, chatId, messages.ClientExerciseRecordNotFoundMessage(params.UserExerciseRecordId))
				return
			}

			newCtx = utils_context.GetContextWithUserExerciseRecord(newCtx, record)
		}

		if params.Limit != 0 {
			newCtx = utils_context.GetContextWithLimit(newCtx, params.Limit)
		}

		next(utils_context.GetContextWithOffset(newCtx, params.Offset), b, update)
	}
}

func (bot *bot) timeoutMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		timeoutDuration := bot.config.RequestTimeout()
		chatId := utils_context.GetChatIdFromContext(ctx)

		childCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		doneCh := make(chan struct{})

		go func() {
			next(childCtx, b, update)
			close(doneCh)
		}()

		select {
		case <-childCtx.Done():
			bot_utils.SendMessage(ctx, b, chatId, messages.RequestTimeoutMessage())
			return
		case <-doneCh:
			return
		}
	}
}

func (bot *bot) panicRecoveryMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
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
					Text:      messages.ErrorMessage(),
					ParseMode: tg_models.ParseModeMarkdown,
				})
			}
		}()

		next(ctx, b, update)
	}
}

func (bot *bot) adminMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{
		bot.answerCallbackQueryMiddleware,
		bot.isRegisteredMiddleware,
		bot.isAdminMiddleware,
		bot.parseParamsMiddleware,
		bot.validateParamsMiddleware,
	}
}

func (bot *bot) userMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{
		bot.timeoutMiddleware,
		bot.answerCallbackQueryMiddleware,
		bot.isRegisteredMiddleware,
		bot.isApprovedMiddleware,
		bot.parseParamsMiddleware,
		bot.validateParamsMiddleware,
	}
}

func (bot *bot) mainMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{
		bot.timeoutMiddleware,
		bot.answerCallbackQueryMiddleware,
		bot.isRegisteredMiddleware,
	}
}

func (bot *bot) commandMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{
		bot.timeoutMiddleware,
	}
}
