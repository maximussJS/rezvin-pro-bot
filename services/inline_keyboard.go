package services

import (
	"fmt"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/models"
)

type IInlineKeyboardService interface {
	AdminMain() *tg_models.InlineKeyboardMarkup
	ProgramMenu() *tg_models.InlineKeyboardMarkup
	UserRegister() *tg_models.InlineKeyboardMarkup
	UserMenu() *tg_models.InlineKeyboardMarkup
	PendingUsersList(users []models.User) *tg_models.InlineKeyboardMarkup
	PendingUserDecide(users models.User) *tg_models.InlineKeyboardMarkup
	ProgramList(programs []models.Program) *tg_models.InlineKeyboardMarkup
	ProgramSelectedMenu(programId uint) *tg_models.InlineKeyboardMarkup
	ProgramExerciseDeleteList(programId uint, exercises []models.Exercise) *tg_models.InlineKeyboardMarkup
}

type inlineKeyboardService struct{}

func NewInlineKeyboardService() *inlineKeyboardService {
	return &inlineKeyboardService{}
}

func (s *inlineKeyboardService) AdminMain() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📖 Програми", CallbackData: callback_data.ProgramMenu},
			},
			{
				{Text: "⏳ Підтвердження клієнтів", CallbackData: callback_data.PendingUsersList},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.BackToStart},
			},
		},
	}
}

func (s *inlineKeyboardService) ProgramMenu() *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список програм", CallbackData: callback_data.ProgramList},
			},
			{
				{Text: "➕ Створити програму", CallbackData: callback_data.ProgramAdd},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.BackToMain},
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
				{Text: "🔙 Назад", CallbackData: callback_data.BackToStart},
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
				{Text: "🔙 Назад", CallbackData: callback_data.BackToStart},
			},
		},
	}
}

func (s *inlineKeyboardService) ProgramList(programs []models.Program) *tg_models.InlineKeyboardMarkup {
	programKb := make([][]tg_models.InlineKeyboardButton, 0, len(programs))

	for _, program := range programs {
		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name,
				CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramSelected, program.Id),
			},
		})
	}

	programKb = append(programKb, []tg_models.InlineKeyboardButton{
		{Text: "🔙 Назад", CallbackData: callback_data.BackToProgramMenu},
	})

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: programKb,
	}
}

func (s *inlineKeyboardService) ProgramSelectedMenu(programId uint) *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список вправ", CallbackData: fmt.Sprintf("%s:%d", callback_data.ExerciseList, programId)},
			},
			{
				{Text: "➕ Додати вправу", CallbackData: fmt.Sprintf("%s:%d", callback_data.ExerciseAdd, programId)},
			},
			{
				{Text: "➖ Видалити вправу", CallbackData: fmt.Sprintf("%s:%d", callback_data.ExerciseDelete, programId)},
			},
			{
				{Text: "📝 Перейменувати програму", CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramRename, programId)},
			},
			{
				{Text: "❌ Видалити програму", CallbackData: fmt.Sprintf("%s:%d", callback_data.ProgramDelete, programId)},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.BackToProgramList},
			},
		},
	}
}

func (s *inlineKeyboardService) ProgramExerciseDeleteList(programId uint, exercises []models.Exercise) *tg_models.InlineKeyboardMarkup {
	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, len(exercises))

	for _, exercise := range exercises {
		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         exercise.Name,
				CallbackData: fmt.Sprintf("%s:%d:%d", callback_data.ExerciseDeleteItem, programId, exercise.Id),
			},
		})
	}

	exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
		{
			Text:         "🔙 Назад",
			CallbackData: callback_data.ProgramBack,
		},
	})

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: exerciseKb,
	}
}

func (s *inlineKeyboardService) PendingUsersList(users []models.User) *tg_models.InlineKeyboardMarkup {
	userKb := make([][]tg_models.InlineKeyboardButton, 0, len(users))

	for _, user := range users {
		text := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

		if user.Username != "" {
			text += fmt.Sprintf(" @%s", user.Username)
		}

		userKb = append(userKb, []tg_models.InlineKeyboardButton{
			{
				Text:         text,
				CallbackData: fmt.Sprintf("%s:%d", callback_data.PendingUsersSelected, user.Id),
			},
		})
	}

	userKb = append(userKb, []tg_models.InlineKeyboardButton{
		{Text: "🔙 Назад", CallbackData: callback_data.BackToMain},
	})

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: userKb,
	}
}

func (s *inlineKeyboardService) PendingUserDecide(users models.User) *tg_models.InlineKeyboardMarkup {
	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "✅ Підтвердити", CallbackData: fmt.Sprintf("%s:%d", callback_data.PendingUsersApprove, users.Id)},
			},
			{
				{Text: "❌ Відхилити", CallbackData: fmt.Sprintf("%s:%d", callback_data.PendingUsersDecline, users.Id)},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.BackToPendingUsersList},
			},
		},
	}
}
