package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func ClientList(clients []models.User, totalClientCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	clientsLen := len(clients)
	clientKb := make([][]tg_models.InlineKeyboardButton, 0, clientsLen)

	for _, client := range clients {
		params := types.NewEmptyParams()

		params.UserId = client.Id

		clientKb = append(clientKb, []tg_models.InlineKeyboardButton{
			{
				Text:         client.GetPrivateName(),
				CallbackData: bot_utils.AddParamsToQueryString(constants.ClientSelected, params),
			},
		})
	}

	clientKb = append(clientKb, GetPaginationButtons(
		clientsLen,
		totalClientCount,
		constants.ClientList,
		limit,
		offset,
		types.NewEmptyParams(),
		types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(clientKb, GetBackButton(constants.MainBackToMain, types.NewEmptyParams())),
	}
}

func ClientSelectedMenu(clientId int64) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Програми клієнта", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientProgramList, params)},
			},
			{
				{Text: "📏 Заміри клієнта", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientMeasureList, params)},
			},
			{
				{Text: "➕ Призначити програму для клієнта", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientProgramAdd, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.BackToClientList},
			},
		},
	}
}

func ClientSelectedOk(clientId int64) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()
	params.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.ClientSelected, params),
		},
	}
}
