package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func ProgramMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список програм", CallbackData: constants.ProgramList},
			},
			{
				{Text: "➕ Створити програму", CallbackData: constants.ProgramAdd},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.MainBackToMain},
			},
		},
	}
}

func ProgramList(programs []models.Program, totalProgramCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	programsLen := len(programs)
	programKb := make([][]tg_models.InlineKeyboardButton, 0, programsLen)

	for _, program := range programs {
		params := types.NewEmptyParams()

		params.ProgramId = program.Id
		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name,
				CallbackData: bot_utils.AddParamsToQueryString(constants.ProgramSelected, params),
			},
		})
	}

	programKb = append(programKb, GetPaginationButtons(
		programsLen,
		totalProgramCount,
		constants.ProgramList,
		limit,
		offset,
		types.NewEmptyParams(),
		types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(constants.ProgramMenu, types.NewEmptyParams())),
	}
}

func ProgramSelectedMenu(programId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.ProgramId = programId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список вправ", CallbackData: bot_utils.AddParamsToQueryString(constants.ExerciseList, params)},
			},
			{
				{Text: "➕ Додати вправу", CallbackData: bot_utils.AddParamsToQueryString(constants.ExerciseAdd, params)},
			},
			{
				{Text: "➖ Видалити вправу", CallbackData: bot_utils.AddParamsToQueryString(constants.ExerciseDelete, params)},
			},
			{
				{Text: "📝 Перейменувати програму", CallbackData: bot_utils.AddParamsToQueryString(constants.ProgramRename, params)},
			},
			{
				{Text: "❌ Видалити програму", CallbackData: bot_utils.AddParamsToQueryString(constants.ProgramDelete, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.BackToProgramList},
			},
		},
	}
}

func ProgramOk(programId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.ProgramId = programId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.ProgramSelected, params),
		},
	}
}

func ProgramDeleteOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.BackToProgramList, types.NewEmptyParams()),
		},
	}
}

func ProgramMenuOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.ProgramMenu, types.NewEmptyParams()),
		},
	}
}
