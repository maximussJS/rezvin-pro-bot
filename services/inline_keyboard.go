package services

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/types/bot"
	bot_utils "rezvin-pro-bot/utils/bot"
)

type IInlineKeyboardService interface {
	AdminMain() *tg_models.InlineKeyboardMarkup
	ProgramMenu() *tg_models.InlineKeyboardMarkup
	UserRegister() *tg_models.InlineKeyboardMarkup
	UserMenu() *tg_models.InlineKeyboardMarkup
	PendingUsersList(users []models.User) *tg_models.InlineKeyboardMarkup
	PendingUserDecide(users models.User) *tg_models.InlineKeyboardMarkup
	ProgramList(programs []models.Program) *tg_models.InlineKeyboardMarkup
	ClientList(clients []models.User) *tg_models.InlineKeyboardMarkup
	ClientSelectedMenu(clientId int64) *tg_models.InlineKeyboardMarkup
	ClientProgramList(clientId int64, programs []models.UserProgram) *tg_models.InlineKeyboardMarkup
	ProgramForClientList(clientId int64, programs []models.Program) *tg_models.InlineKeyboardMarkup
	ClientSelectedProgramMenu(clientId int64, program models.UserProgram) *tg_models.InlineKeyboardMarkup
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
				{Text: "🏋️ Клієнти", CallbackData: callback_data.ClientList},
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
		params := bot_types.NewEmptyParams()

		params.ProgramId = program.Id
		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ProgramSelected, params),
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
	params := bot_types.NewEmptyParams()

	params.ProgramId = programId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Список вправ", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ExerciseList, params)},
			},
			{
				{Text: "➕ Додати вправу", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ExerciseAdd, params)},
			},
			{
				{Text: "➖ Видалити вправу", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ExerciseDelete, params)},
			},
			{
				{Text: "📝 Перейменувати програму", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ProgramRename, params)},
			},
			{
				{Text: "❌ Видалити програму", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ProgramDelete, params)},
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
		params := bot_types.NewEmptyParams()

		params.ProgramId = programId
		params.ExerciseId = exercise.Id

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         exercise.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ExerciseDeleteItem, params),
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
		params := bot_types.NewEmptyParams()

		params.UserId = user.Id

		userKb = append(userKb, []tg_models.InlineKeyboardButton{
			{
				Text:         user.GetPrivateName(),
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.PendingUsersSelected, params),
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
	params := bot_types.NewEmptyParams()

	params.UserId = users.Id

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

func (s *inlineKeyboardService) ClientList(clients []models.User) *tg_models.InlineKeyboardMarkup {
	clientKb := make([][]tg_models.InlineKeyboardButton, 0, len(clients))

	for _, client := range clients {
		params := bot_types.NewEmptyParams()

		params.UserId = client.Id

		clientKb = append(clientKb, []tg_models.InlineKeyboardButton{
			{
				Text:         client.GetPrivateName(),
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientSelected, params),
			},
		})
	}

	clientKb = append(clientKb, []tg_models.InlineKeyboardButton{
		{Text: "🔙 Назад", CallbackData: callback_data.BackToMain},
	})

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: clientKb,
	}
}

func (s *inlineKeyboardService) ClientSelectedMenu(clientId int64) *tg_models.InlineKeyboardMarkup {
	params := bot_types.NewEmptyParams()

	params.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "📋 Дивитись програми клієнта", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramList, params)},
			},
			{
				{Text: "➕ Додати програму для клієнта", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramAdd, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: callback_data.BackToClientList},
			},
		},
	}
}

func (s *inlineKeyboardService) ClientSelectedProgramMenu(clientId int64, program models.UserProgram) *tg_models.InlineKeyboardMarkup {
	params := bot_types.NewEmptyParams()

	params.UserId = clientId
	params.UserProgramId = program.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultList, params)},
			},
			{
				{Text: "✍️ Внести результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultModify, params)},
			},
			{
				{Text: "➖ Видалити програму", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramDelete, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramList, params)},
			},
		},
	}
}

func (s *inlineKeyboardService) ClientProgramList(clientId int64, programs []models.UserProgram) *tg_models.InlineKeyboardMarkup {
	programKb := make([][]tg_models.InlineKeyboardButton, 0, len(programs))

	for _, program := range programs {
		params := bot_types.NewEmptyParams()

		params.UserId = clientId
		params.UserProgramId = program.Id

		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Program.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramSelected, params),
			},
		})
	}

	params := bot_types.NewEmptyParams()

	params.UserId = clientId

	programKb = append(programKb, []tg_models.InlineKeyboardButton{
		{Text: "🔙 Назад", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientSelected, params)},
	})

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: programKb,
	}
}

func (s *inlineKeyboardService) ProgramForClientList(clientId int64, programs []models.Program) *tg_models.InlineKeyboardMarkup {
	programKb := make([][]tg_models.InlineKeyboardButton, 0, len(programs))

	for _, program := range programs {
		params := bot_types.NewEmptyParams()

		params.UserId = clientId
		params.ProgramId = program.Id
		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramAssign, params),
			},
		})
	}

	params := bot_types.NewEmptyParams()

	params.UserId = clientId

	programKb = append(programKb, []tg_models.InlineKeyboardButton{
		{Text: "🔙 Назад", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientSelected, params)},
	})

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: programKb,
	}
}
