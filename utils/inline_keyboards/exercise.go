package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/models"
	bot_types "rezvin-pro-bot/types/bot"
	bot_utils "rezvin-pro-bot/utils/bot"
)

func ProgramExerciseDeleteList(programId uint, exercises []models.Exercise, totalExerciseCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	exercisesLen := len(exercises)
	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, exercisesLen)

	for _, exercise := range exercises {
		params := bot_types.NewEmptyParams()

		params.ProgramId = programId
		params.ExerciseId = exercise.Id

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         exercise.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ExerciseDeleteItem, params),
			},
		})
	}

	nextParams := bot_types.NewEmptyParams()
	nextParams.ProgramId = programId

	previousParams := bot_types.NewEmptyParams()
	previousParams.ProgramId = programId

	exerciseKb = append(exerciseKb, GetPaginationButtons(
		exercisesLen,
		totalExerciseCount,
		callback_data.ExerciseDelete,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(exerciseKb, GetBackButton(callback_data.ProgramBack, bot_types.NewEmptyParams())),
	}
}
