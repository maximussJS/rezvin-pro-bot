package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func PendingUsersList(users []models.User, totalUsersCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	usersLen := len(users)
	userKb := make([][]tg_models.InlineKeyboardButton, 0, usersLen)

	for _, user := range users {
		params := types.NewEmptyParams()

		params.UserId = user.Id

		userKb = append(userKb, []tg_models.InlineKeyboardButton{
			{
				Text:         user.GetPrivateName(),
				CallbackData: bot_utils.AddParamsToQueryString(constants.PendingUsersSelected, params),
			},
		})
	}

	userKb = append(userKb, GetPaginationButtons(
		usersLen,
		totalUsersCount,
		constants.PendingUsersList,
		limit,
		offset,
		types.NewEmptyParams(),
		types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(userKb, GetBackButton(constants.MainBackToMain, types.NewEmptyParams())),
	}
}

func PendingUserDecide(user models.User) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.UserId = user.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "✅ Підтвердити", CallbackData: bot_utils.AddParamsToQueryString(constants.PendingUsersApprove, params)},
			},
			{
				{Text: "❌ Відхилити", CallbackData: bot_utils.AddParamsToQueryString(constants.PendingUsersDecline, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.BackToPendingUsersList},
			},
		},
	}
}

func PendingUsersOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.PendingUsersList, types.NewEmptyParams()),
		},
	}
}
