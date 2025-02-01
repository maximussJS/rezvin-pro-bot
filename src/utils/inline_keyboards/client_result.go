package inline_keyboards

import (
	"fmt"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func ClientResultsExerciseList(clientId int64, userProgramId uint, exercises []models.Exercise, totalExerciseCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	exercisesLen := len(exercises)
	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, exercisesLen)

	for _, exercise := range exercises {
		params := types.NewEmptyParams()

		params.UserId = clientId
		params.UserProgramId = userProgramId
		params.ExerciseId = exercise.Id

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         exercise.Name,
				CallbackData: bot_utils.AddParamsToQueryString(constants.ClientResultExerciseSelected, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()
	nextParams.UserId = clientId
	nextParams.UserProgramId = userProgramId

	previousParams := types.NewEmptyParams()
	previousParams.UserId = clientId
	previousParams.UserProgramId = userProgramId

	exerciseKb = append(exerciseKb, GetPaginationButtons(
		exercisesLen,
		totalExerciseCount,
		constants.ClientResultExercisesList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(exerciseKb, GetBackButton(constants.ClientSelected, backParams)),
	}
}

func ClientResultExerciseSelectedOk(clientId int64, records []models.UserResult) *tg_models.InlineKeyboardMarkup {
	recordsLen := len(records)
	recordsKb := make([][]tg_models.InlineKeyboardButton, 0, recordsLen)

	for _, record := range records {
		params := types.NewEmptyParams()

		params.UserId = clientId
		params.UserProgramId = record.UserProgramId
		params.ExerciseId = record.ExerciseId
		params.UserResultId = record.Id

		recordsKb = append(recordsKb, []tg_models.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d повторень", record.Reps),
				CallbackData: bot_utils.AddParamsToQueryString(constants.ClientResultExerciseReps, params),
			},
		})
	}

	backParams := types.NewEmptyParams()
	backParams.UserId = clientId
	backParams.UserProgramId = records[0].UserProgramId
	backParams.ExerciseId = records[0].ExerciseId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(recordsKb, GetBackButton(constants.ClientResultExercisesList, backParams)),
	}
}
