package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	bot_types "rezvin-pro-bot/types/bot"
	bot_utils "rezvin-pro-bot/utils/bot"
)

func GetPaginationButtons(
	itemsLength int, totalCount int64,
	callBackData string,
	limit, offset int,
	nextParams, previousParams *bot_types.Params,
) []tg_models.InlineKeyboardButton {
	kb := make([]tg_models.InlineKeyboardButton, 0, 2)

	if offset >= limit {
		previousParams.Offset = offset - limit
		previousParams.Limit = limit

		kb = append(kb, tg_models.InlineKeyboardButton{
			Text:         "⬅️ Попередні",
			CallbackData: bot_utils.AddParamsToQueryString(callBackData, previousParams),
		})
	}

	if itemsLength == limit && int64(offset) < totalCount-int64(limit) {
		nextParams.Offset = offset + limit
		nextParams.Limit = limit

		kb = append(kb, tg_models.InlineKeyboardButton{
			Text:         "➡️ Наступні",
			CallbackData: bot_utils.AddParamsToQueryString(callBackData, nextParams),
		})
	}

	return kb
}

func GetBackButton(callBackData string, params *bot_types.Params) []tg_models.InlineKeyboardButton {
	return []tg_models.InlineKeyboardButton{
		{Text: "🔙 Назад", CallbackData: bot_utils.AddParamsToQueryString(callBackData, params)},
	}
}

func GetOkButton(callBackData string, params *bot_types.Params) []tg_models.InlineKeyboardButton {
	return []tg_models.InlineKeyboardButton{
		{Text: "✅ Ок", CallbackData: bot_utils.AddParamsToQueryString(callBackData, params)},
	}
}
