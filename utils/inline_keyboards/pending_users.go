package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/models"
	bot_types "rezvin-pro-bot/types/bot"
	bot_utils "rezvin-pro-bot/utils/bot"
)

func PendingUsersList(users []models.User, totalUsersCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	usersLen := len(users)
	userKb := make([][]tg_models.InlineKeyboardButton, 0, usersLen)

	for _, user := range users {
		params := bot_types.NewEmptyParams()

		params.UserId = user.Id

		userKb = append(userKb, []tg_models.InlineKeyboardButton{
			{
				Text:         user.GetPrivateName(),
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.PendingUsersSelected, params),
			},
		})
	}

	userKb = append(userKb, GetPaginationButtons(
		usersLen,
		totalUsersCount,
		callback_data.PendingUsersList,
		limit,
		offset,
		bot_types.NewEmptyParams(),
		bot_types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(userKb, GetBackButton(callback_data.MainBackToMain, bot_types.NewEmptyParams())),
	}
}

func PendingUserDecide(user models.User) *tg_models.InlineKeyboardMarkup {
	params := bot_types.NewEmptyParams()

	params.UserId = user.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "✅ Підтвердити", CallbackData: bot_utils.AddParamsToQueryString(callback_data.PendingUsersApprove, params)},
			},
			{
				{Text: "❌ Відхилити", CallbackData: bot_utils.AddParamsToQueryString(callback_data.PendingUsersDecline, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.BackToPendingUsersList},
			},
		},
	}
}

func PendingUsersOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(callback_data.PendingUsersList, bot_types.NewEmptyParams()),
		},
	}
}
