package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func ClientMeasuresList(clientId int64, measures []models.Measure, totalMeasureCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	measuresLen := len(measures)
	measuresKb := make([][]tg_models.InlineKeyboardButton, 0, measuresLen)

	for _, measure := range measures {
		params := types.NewEmptyParams()

		params.UserId = clientId
		params.MeasureId = measure.Id

		measuresKb = append(measuresKb, []tg_models.InlineKeyboardButton{
			{
				Text:         measure.Name,
				CallbackData: bot_utils.AddParamsToQueryString(constants.ClientMeasureSelected, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()
	nextParams.UserId = clientId

	previousParams := types.NewEmptyParams()
	previousParams.UserId = clientId

	measuresKb = append(measuresKb, GetPaginationButtons(
		measuresLen,
		totalMeasureCount,
		constants.ClientMeasureList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(measuresKb, GetBackButton(constants.ClientSelected, backParams)),
	}
}

func ClientMeasureMenu(clientId int64, measureId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.UserId = clientId
	params.MeasureId = measureId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати заміру", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientMeasureResult, params)},
			},
			{
				{Text: "➕️ Внести результат заміру", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientMeasureAdd, params)},
			},
			{
				{Text: "➖ Видалити останній результат заміру", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientMeasureDelete, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientMeasureList, params)},
			},
		},
	}
}

func ClientMeasureOk(clientId int64, measureId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.UserId = clientId
	params.MeasureId = measureId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.ClientMeasureSelected, params),
		},
	}
}
