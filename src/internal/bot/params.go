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

		if params.UserId != 0 {
			user := bot.userRepository.GetById(ctx, params.UserId)

			if user == nil {
				msg := messages.UserNotFoundMessage(params.UserId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			ctx = utils_context.GetContextWithUser(ctx, user)
		}

		if params.ProgramId != 0 {
			program := bot.programRepository.GetById(ctx, params.ProgramId)

			if program == nil {
				msg := messages.ProgramNotFoundMessage(params.ProgramId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			ctx = utils_context.GetContextWithProgram(ctx, program)
		}

		if params.ExerciseId != 0 {
			exercise := bot.exerciseRepository.GetById(ctx, params.ExerciseId)

			if exercise == nil {
				msg := messages.ExerciseNotFoundMessage(params.ExerciseId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			ctx = utils_context.GetContextWithExercise(ctx, exercise)
		}

		if params.UserProgramId != 0 {
			userProgram := bot.userProgramRepository.GetById(ctx, params.UserProgramId)

			if userProgram == nil {
				msg := messages.ClientProgramNotFoundMessage(params.UserProgramId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			ctx = utils_context.GetContextWithUserProgram(ctx, userProgram)
		}

		if params.UserResultId != 0 {
			record := bot.userResultRepository.GetById(ctx, params.UserResultId)

			if record == nil {
				msg := messages.ClientResultNotFoundMessage(params.UserResultId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			ctx = utils_context.GetContextWithUserResult(ctx, record)
		}

		if params.MeasureId != 0 {
			measure := bot.measureRepository.GetById(ctx, params.MeasureId)

			if measure == nil {
				msg := messages.MeasureNotFoundMessage(params.MeasureId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			ctx = utils_context.GetContextWithMeasure(ctx, measure)
		}

		if params.UserMeasureId != 0 {
			userMeasure := bot.userMeasureRepository.GetById(ctx, params.UserMeasureId)

			if userMeasure == nil {
				msg := messages.ClientMeasureNotFoundMessage(params.UserMeasureId)
				kb := inline_keyboards.StartOk()

				bot.senderService.SendWithKb(ctx, b, chatId, msg, kb)
				return
			}

			ctx = utils_context.GetContextWithUserMeasure(ctx, userMeasure)
		}

		if params.Reps != constants.Zero {
			ctx = utils_context.GetContextWithReps(ctx, params.Reps)
		}

		if params.Limit != 0 {
			ctx = utils_context.GetContextWithLimit(ctx, params.Limit)
		}

		next(utils_context.GetContextWithOffset(ctx, params.Offset), b, update)
	}
}
