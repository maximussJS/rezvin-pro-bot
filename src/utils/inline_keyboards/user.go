package inline_keyboards

import (
	"fmt"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants/callback_data"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func UserProgramList(programs []models.UserProgram, totalProgramCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	programsLen := len(programs)

	programKb := make([][]tg_models.InlineKeyboardButton, 0, programsLen)

	for _, program := range programs {
		params := types.NewEmptyParams()

		params.UserProgramId = program.Id

		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name(),
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.UserProgramSelected, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()
	nextParams.UserProgramId = programs[0].Id
	previousParams := types.NewEmptyParams()
	previousParams.UserProgramId = programs[0].Id

	programKb = append(programKb, GetPaginationButtons(
		programsLen,
		totalProgramCount,
		callback_data.UserProgramList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(callback_data.MainBackToMain, types.NewEmptyParams())),
	}
}

func UserProgramMenu(userProgram models.UserProgram) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.UserProgramId = userProgram.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.UserResultList, params)},
			},
			{
				{Text: "✍️ Внести результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.UserResultModifyExerciseList, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.MainBackToMain},
			},
		},
	}
}

func UserMenuOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(callback_data.MainBackToMain, types.NewEmptyParams()),
		},
	}
}

func UserProgramMenuOk(userProgramId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()
	params.UserProgramId = userProgramId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(callback_data.UserProgramSelected, params),
		},
	}
}

func UserProgramListOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(callback_data.UserProgramList, types.NewEmptyParams()),
		},
	}
}

func UserProgramResultsModifyExerciseList(userProgramId uint, exercises []models.Exercise, totalExerciseCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	exercisesLen := len(exercises)
	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, exercisesLen)

	for _, exercise := range exercises {
		params := types.NewEmptyParams()

		params.UserProgramId = userProgramId
		params.ExerciseId = exercise.Id

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         exercise.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.UserResultModifyExerciseSelected, params),
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
		callback_data.UserResultModifyExerciseList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()
	backParams.UserProgramId = userProgramId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(exerciseKb, GetBackButton(callback_data.UserProgramSelected, backParams)),
	}
}

func UserProgramResultModifyExerciseSelectedOk(records []models.UserExerciseRecord) *tg_models.InlineKeyboardMarkup {
	recordsLen := len(records)
	recordsKb := make([][]tg_models.InlineKeyboardButton, 0, recordsLen)

	for _, record := range records {
		params := types.NewEmptyParams()

		params.UserProgramId = record.UserProgramId
		params.ExerciseId = record.ExerciseId
		params.UserExerciseRecordId = record.Id

		recordsKb = append(recordsKb, []tg_models.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d повторень", record.Reps),
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.UserResultModifyExerciseRepsModify, params),
			},
		})
	}

	backParams := types.NewEmptyParams()
	backParams.UserProgramId = records[0].UserProgramId
	backParams.ExerciseId = records[0].ExerciseId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(recordsKb, GetBackButton(callback_data.UserResultModifyExerciseList, backParams)),
	}
}
