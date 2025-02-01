package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func MeasureMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список замірів", CallbackData: constants.MeasureList},
			},
			{
				{Text: "➕ Створити замір", CallbackData: constants.MeasureAdd},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.MainBackToMain},
			},
		},
	}
}

func MeasureOk(measureId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.MeasureId = measureId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.MeasureSelected, params),
		},
	}
}

func MeasureMenuOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.MeasureMenu, types.NewEmptyParams()),
		},
	}
}

func MeasureDeleteOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.BackToMeasureList, types.NewEmptyParams()),
		},
	}
}

func MeasureList(measures []models.Measure, totalMeasureCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	measuresLen := len(measures)
	measuresKb := make([][]tg_models.InlineKeyboardButton, 0, measuresLen)

	for _, measure := range measures {
		params := types.NewEmptyParams()

		params.MeasureId = measure.Id
		measuresKb = append(measuresKb, []tg_models.InlineKeyboardButton{
			{
				Text:         measure.Name,
				CallbackData: bot_utils.AddParamsToQueryString(constants.MeasureSelected, params),
			},
		})
	}

	measuresKb = append(measuresKb, GetPaginationButtons(
		measuresLen,
		totalMeasureCount,
		constants.MeasureList,
		limit,
		offset,
		types.NewEmptyParams(),
		types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(measuresKb, GetBackButton(constants.MeasureMenu, types.NewEmptyParams())),
	}
}

func MeasureSelectedMenu(measureId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.MeasureId = measureId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "⏱️ Перейменувати замір", CallbackData: bot_utils.AddParamsToQueryString(constants.MeasureRename, params)},
			},
			{
				{Text: "📏 Змінити одиниці виміру", CallbackData: bot_utils.AddParamsToQueryString(constants.MeasureChangeUnits, params)},
			},
			{
				{Text: "❌ Видалити замір", CallbackData: bot_utils.AddParamsToQueryString(constants.MeasureDelete, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.BackToMeasureList},
			},
		},
	}
}
