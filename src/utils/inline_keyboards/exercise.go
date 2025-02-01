package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func ExerciseDeleteList(programId uint, exercises []models.Exercise, totalExerciseCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	exercisesLen := len(exercises)
	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, exercisesLen)

	for _, exercise := range exercises {
		params := types.NewEmptyParams()

		params.ProgramId = programId
		params.ExerciseId = exercise.Id

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         exercise.Name,
				CallbackData: bot_utils.AddParamsToQueryString(constants.ExerciseDeleteItem, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()
	nextParams.ProgramId = programId

	previousParams := types.NewEmptyParams()
	previousParams.ProgramId = programId

	exerciseKb = append(exerciseKb, GetPaginationButtons(
		exercisesLen,
		totalExerciseCount,
		constants.ExerciseDelete,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()
	backParams.ProgramId = programId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(exerciseKb, GetBackButton(constants.ProgramSelected, backParams)),
	}
}

func ExerciseOk(programId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.ProgramId = programId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.ProgramSelected, params),
		},
	}
}
