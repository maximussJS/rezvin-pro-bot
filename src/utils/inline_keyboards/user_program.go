package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
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
				CallbackData: bot_utils.AddParamsToQueryString(constants.UserProgramSelected, params),
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
		constants.UserProgramList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(constants.MainBackToMain, types.NewEmptyParams())),
	}
}

func UserProgramMenu(userProgram models.UserProgram) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.UserProgramId = userProgram.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати", CallbackData: bot_utils.AddParamsToQueryString(constants.UserResultList, params)},
			},
			{
				{Text: "✍️ Внести результати", CallbackData: bot_utils.AddParamsToQueryString(constants.UserResultExerciseList, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.MainBackToMain},
			},
		},
	}
}

func UserProgramMenuOk(userProgramId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()
	params.UserProgramId = userProgramId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.UserProgramSelected, params),
		},
	}
}

func UserProgramListOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.UserProgramList, types.NewEmptyParams()),
		},
	}
}
