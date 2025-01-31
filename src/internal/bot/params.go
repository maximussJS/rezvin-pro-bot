package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	bot_utils "rezvin-pro-bot/src/utils/bot"
	utils_context "rezvin-pro-bot/src/utils/context"
	"rezvin-pro-bot/src/utils/inline_keyboards"
	"rezvin-pro-bot/src/utils/messages"
)

func (bot *bot) parseParamsMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
		callbackQueryData := update.CallbackQuery.Data
		chatId := utils_context.GetChatIdFromContext(ctx)

		params, err := bot_utils.ParseParamsFromQueryString(callbackQueryData)

		if err != nil {
			bot.logger.Error(fmt.Sprintf("Failed to parse params: %s", callbackQueryData))
			msg := messages.ParamsErrorMessage(err)
			kb := inline_keyboards.StartOk()

			bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
			return
		}

		if params == nil {
			bot.logger.Error(fmt.Sprintf("Failed to parse params: %s", update.Message.Text))
			msg := messages.ErrorMessage()
			kb := inline_keyboards.StartOk()

			bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
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
				msg := messages.UserNotFoundMessage(params.UserId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			newCtx = utils_context.GetContextWithUser(newCtx, user)
		}

		if params.ProgramId != 0 {
			program := bot.programRepository.GetById(ctx, params.ProgramId)

			if program == nil {
				msg := messages.ProgramNotFoundMessage(params.ProgramId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			newCtx = utils_context.GetContextWithProgram(newCtx, program)
		}

		if params.ExerciseId != 0 {
			exercise := bot.exerciseRepository.GetById(ctx, params.ExerciseId)

			if exercise == nil {
				msg := messages.ExerciseNotFoundMessage(params.ExerciseId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			newCtx = utils_context.GetContextWithExercise(newCtx, exercise)
		}

		if params.UserProgramId != 0 {
			userProgram := bot.userProgramRepository.GetById(ctx, params.UserProgramId)

			if userProgram == nil {
				msg := messages.ClientProgramNotFoundMessage(params.UserProgramId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			newCtx = utils_context.GetContextWithUserProgram(newCtx, userProgram)
		}

		if params.UserExerciseRecordId != 0 {
			record := bot.userExerciseRecordRepository.GetById(ctx, params.UserExerciseRecordId)

			if record == nil {
				msg := messages.ClientExerciseRecordNotFoundMessage(params.UserExerciseRecordId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			newCtx = utils_context.GetContextWithUserExerciseRecord(newCtx, record)
		}

		if params.Reps != constants.Zero {
			newCtx = utils_context.GetContextWithReps(newCtx, params.Reps)
		}

		if params.Limit != 0 {
			newCtx = utils_context.GetContextWithLimit(newCtx, params.Limit)
		}

		next(utils_context.GetContextWithOffset(newCtx, params.Offset), b, update)
	}
}
