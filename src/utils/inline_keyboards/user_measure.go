package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func UserMeasureMenu(measureId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.MeasureId = measureId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати заміру", CallbackData: bot_utils.AddParamsToQueryString(constants.UserMeasureResult, params)},
			},
			{
				{Text: "➕️ Внести результат заміру", CallbackData: bot_utils.AddParamsToQueryString(constants.UserMeasureAdd, params)},
			},
			{
				{Text: "➖ Видалити останній результат заміру", CallbackData: bot_utils.AddParamsToQueryString(constants.UserMeasureDelete, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: bot_utils.AddParamsToQueryString(constants.UserMeasureList, params)},
			},
		},
	}
}

func UserMeasureOk(measureId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.MeasureId = measureId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.UserMeasureSelected, params),
		},
	}
}

func UserMeasuresList(measures []models.Measure, totalMeasureCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	measuresLen := len(measures)
	measuresKb := make([][]tg_models.InlineKeyboardButton, 0, measuresLen)

	for _, measure := range measures {
		params := types.NewEmptyParams()

		params.MeasureId = measure.Id

		measuresKb = append(measuresKb, []tg_models.InlineKeyboardButton{
			{
				Text:         measure.Name,
				CallbackData: bot_utils.AddParamsToQueryString(constants.UserMeasureSelected, params),
			},
		})
	}

	measuresKb = append(measuresKb, GetPaginationButtons(
		measuresLen,
		totalMeasureCount,
		constants.UserMeasureList,
		limit,
		offset,
		types.NewEmptyParams(),
		types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(measuresKb, GetBackButton(constants.MainBackToMain, types.NewEmptyParams())),
	}
}
