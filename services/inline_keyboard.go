package services

import (
	"fmt"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/models"
)

type IInlineKeyboardService interface {
	AdminMenu() *tg_models.InlineKeyboardMarkup
	AdminProgramMenu() *tg_models.InlineKeyboardMarkup
	UserRegister() *tg_models.InlineKeyboardMarkup
	UserMenu() *tg_models.InlineKeyboardMarkup
	ProgramMenu(programs []models.Program) *tg_models.InlineKeyboardMarkup
	ProgramSelectedMenu(programId uint) *tg_models.InlineKeyboardMarkup
}

type inlineKeyboardService struct{}

func NewInlineKeyboardService() *inlineKeyboardService {
	return &inlineKeyboardService{}
}

func (s *inlineKeyboardService) AdminMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📖 Програми", CallbackData: callback_data.AdminProgramMenu},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.AdminBack},
			},
		},
	}
}

func (s *inlineKeyboardService) AdminProgramMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список програм", CallbackData: callback_data.AdminProgramMenuList},
			},
			{
				{Text: "➕ Створити програму", CallbackData: callback_data.AdminProgramMenuAdd},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.AdminProgramMenuBack},
			},
		},
	}
}

func (s *inlineKeyboardService) UserRegister() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📲 Реєстрація", CallbackData: callback_data.UserRegister},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.UserBack},
			},
		},
	}
}

func (s *inlineKeyboardService) UserMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати", CallbackData: callback_data.UserGetResults},
			},
			{
				{Text: "✍️ Внести результати", CallbackData: callback_data.UserAddResults},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.UserBack},
			},
		},
	}
}

func (s *inlineKeyboardService) ProgramMenu(programs []models.Program) *tg_models.InlineKeyboardMarkup {
	programKb := make([][]tg_models.InlineKeyboardButton, 0, len(programs))

	for _, program := range programs {
		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name,
				CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramSelected, program.Id),
			},
		})
	}

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: programKb,
	}
}

func (s *inlineKeyboardService) ProgramSelectedMenu(programId uint) *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список вправ", CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramExerciseList, programId)},
			},
			{
				{Text: "➕ Додати вправу", CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramExerciseAdd, programId)},
			},
			{
				{Text: "➖ Видалити вправу", CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramExerciseDelete, programId)},
			},
			{
				{Text: "📝 Перейменувати програму", CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramRename, programId)},
			},
			{
				{Text: "❌ Видалити програму", CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramDelete, programId)},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.ProgramBack},
			},
		},
	}
}
