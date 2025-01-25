package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants/callback_data"
	"rezvin-pro-bot/src/types"
)

func AdminMain() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📖 Програми", CallbackData: callback_data.ProgramMenu},
			},
			{
				{Text: "⏳ Підтвердження клієнтів", CallbackData: callback_data.PendingUsersList},
			},
			{
				{Text: "🏋️ Клієнти", CallbackData: callback_data.ClientList},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.MainBackToStart},
			},
		},
	}
}

func UserRegister() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📲 Реєстрація", CallbackData: callback_data.UserRegister},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.MainBackToStart},
			},
		},
	}
}

func UserMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Мої програми", CallbackData: callback_data.UserProgramList},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.MainBackToStart},
			},
		},
	}
}

func MainOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(callback_data.MainBackToMain, types.NewEmptyParams()),
		},
	}
}

func StartOk() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(callback_data.MainBackToStart, types.NewEmptyParams()),
		},
	}
}
