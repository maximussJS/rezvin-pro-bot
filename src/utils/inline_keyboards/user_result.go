package inline_keyboards

import (
	"fmt"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func UserResultExerciseList(userProgramId uint, exercises []models.Exercise, totalExerciseCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	exercisesLen := len(exercises)
	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, exercisesLen)

	for _, exercise := range exercises {
		params := types.NewEmptyParams()

		params.UserProgramId = userProgramId
		params.ExerciseId = exercise.Id

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         exercise.Name,
				CallbackData: bot_utils.AddParamsToQueryString(constants.UserResultExerciseSelected, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()
	nextParams.UserProgramId = userProgramId

	previousParams := types.NewEmptyParams()
	previousParams.UserProgramId = userProgramId

	exerciseKb = append(exerciseKb, GetPaginationButtons(
		exercisesLen,
		totalExerciseCount,
		constants.UserResultExerciseList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()
	backParams.UserProgramId = userProgramId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(exerciseKb, GetBackButton(constants.UserProgramSelected, backParams)),
	}
}

func UserResultExerciseSelectedOk(records []models.UserResult) *tg_models.InlineKeyboardMarkup {
	recordsLen := len(records)
	recordsKb := make([][]tg_models.InlineKeyboardButton, 0, recordsLen)

	for _, record := range records {
		params := types.NewEmptyParams()

		params.UserProgramId = record.UserProgramId
		params.ExerciseId = record.ExerciseId
		params.UserResultId = record.Id

		recordsKb = append(recordsKb, []tg_models.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d повторень", record.Reps),
				CallbackData: bot_utils.AddParamsToQueryString(constants.UserResultExerciseReps, params),
			},
		})
	}

	backParams := types.NewEmptyParams()
	backParams.UserProgramId = records[0].UserProgramId
	backParams.ExerciseId = records[0].ExerciseId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(recordsKb, GetBackButton(constants.UserResultExerciseList, backParams)),
	}
}
