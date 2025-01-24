package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/models"
	bot_types "rezvin-pro-bot/types/bot"
	bot_utils "rezvin-pro-bot/utils/bot"
)

func ProgramMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список програм", CallbackData: callback_data.ProgramList},
			},
			{
				{Text: "➕ Створити програму", CallbackData: callback_data.ProgramAdd},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.MainBackToMain},
			},
		},
	}
}

func ProgramList(programs []models.Program, totalProgramCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	programsLen := len(programs)
	programKb := make([][]tg_models.InlineKeyboardButton, 0, programsLen)

	for _, program := range programs {
		params := bot_types.NewEmptyParams()

		params.ProgramId = program.Id
		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ProgramSelected, params),
			},
		})
	}

	programKb = append(programKb, GetPaginationButtons(
		programsLen,
		totalProgramCount,
		callback_data.ProgramList,
		limit,
		offset,
		bot_types.NewEmptyParams(),
		bot_types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(callback_data.MainBackToMain, bot_types.NewEmptyParams())),
	}
}

func ProgramSelectedMenu(programId uint) *tg_models.InlineKeyboardMarkup {
	params := bot_types.NewEmptyParams()

	params.ProgramId = programId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список вправ", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ExerciseList, params)},
			},
			{
				{Text: "➕ Додати вправу", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ExerciseAdd, params)},
			},
			{
				{Text: "➖ Видалити вправу", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ExerciseDelete, params)},
			},
			{
				{Text: "📝 Перейменувати програму", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ProgramRename, params)},
			},
			{
				{Text: "❌ Видалити програму", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ProgramDelete, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.BackToProgramList},
			},
		},
	}
}
