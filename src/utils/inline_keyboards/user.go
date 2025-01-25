package inline_keyboards

import (
	"fmt"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants/callback_data"
	models2 "rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func UserProgramResultsModifyList(exercises []models2.UserExerciseRecord, totalExerciseCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	exercisesLen := len(exercises)

	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, exercisesLen)

	for _, exercise := range exercises {
		params := types.NewEmptyParams()

		params.UserProgramId = exercise.UserProgramId
		params.UserExerciseRecordId = exercise.Id

		text := fmt.Sprintf("%s (%d повторень)", exercise.Name(), exercise.Reps)

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         text,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.UserResultModifySelected, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()

	nextParams.UserProgramId = exercises[0].UserProgramId

	previousParams := types.NewEmptyParams()

	previousParams.UserProgramId = exercises[0].UserProgramId

	exerciseKb = append(exerciseKb, GetPaginationButtons(
		exercisesLen,
		totalExerciseCount,
		callback_data.UserResultModifyList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()

	backParams.UserProgramId = exercises[0].UserProgramId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(exerciseKb, GetBackButton(callback_data.UserProgramSelected, backParams)),
	}
}

func UserProgramList(programs []models2.UserProgram, totalProgramCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
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

func UserProgramMenu(userProgram models2.UserProgram) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.UserProgramId = userProgram.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.UserResultList, params)},
			},
			{
				{Text: "✍️ Внести результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.UserResultModifyList, params)},
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
