package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/types"
)

func AdminMain() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📖 Програми", CallbackData: constants.ProgramMenu},
			},
			{
				{Text: "⏱️ Заміри", CallbackData: constants.MeasureMenu},
			},
			{
				{Text: "⏳ Підтвердження клієнтів", CallbackData: constants.PendingUsersList},
			},
			{
				{Text: "🏋️ Клієнти", CallbackData: constants.ClientList},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.MainBackToStart},
			},
		},
	}
}

func UserRegister() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📲 Реєстрація", CallbackData: constants.UserRegister},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.MainBackToStart},
			},
		},
	}
}

func UserMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Мої програми", CallbackData: constants.UserProgramList},
			},
			{
				{Text: "⏱️ Заміри", CallbackData: constants.UserMeasureList},
			},
			{
				{Text: "🔙 Назад", CallbackData: constants.MainBackToStart},
			},
		},
	}
}

func MainOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.MainBackToMain, types.NewEmptyParams()),
		},
	}
}

func StartOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.MainBackToStart, types.NewEmptyParams()),
		},
	}
}
