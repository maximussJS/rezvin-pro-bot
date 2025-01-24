package inline_keyboards

import (
	"fmt"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/models"
	bot_types "rezvin-pro-bot/types/bot"
	bot_utils "rezvin-pro-bot/utils/bot"
)

func ClientList(clients []models.User, totalClientCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	clientsLen := len(clients)
	clientKb := make([][]tg_models.InlineKeyboardButton, 0, clientsLen)

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

	clientKb = append(clientKb, GetPaginationButtons(
		clientsLen,
		totalClientCount,
		callback_data.ClientList,
		limit,
		offset,
		bot_types.NewEmptyParams(),
		bot_types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(clientKb, GetBackButton(callback_data.MainBackToMain, bot_types.NewEmptyParams())),
	}
}

func ClientSelectedMenu(clientId int64) *tg_models.InlineKeyboardMarkup {
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

func ClientSelectedProgramMenu(clientId int64, program models.UserProgram) *tg_models.InlineKeyboardMarkup {
	params := bot_types.NewEmptyParams()

	params.UserId = clientId
	params.UserProgramId = program.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultList, params)},
			},
			{
				{Text: "✍️ Внести результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultModifyList, params)},
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

func ClientProgramList(clientId int64, programs []models.UserProgram, totalProgramCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	programsLen := len(programs)
	programKb := make([][]tg_models.InlineKeyboardButton, 0, programsLen)

	for _, program := range programs {
		params := bot_types.NewEmptyParams()

		params.UserId = clientId
		params.UserProgramId = program.Id

		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name(),
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramSelected, params),
			},
		})
	}

	nextParams := bot_types.NewEmptyParams()
	nextParams.UserId = clientId

	previousParams := bot_types.NewEmptyParams()
	previousParams.UserId = clientId

	programKb = append(programKb, GetPaginationButtons(
		programsLen,
		totalProgramCount,
		callback_data.ClientProgramList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := bot_types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(callback_data.ClientSelected, backParams)),
	}
}

func ProgramForClientList(clientId int64, programs []models.Program, totalProgramCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	programsLen := len(programs)
	programKb := make([][]tg_models.InlineKeyboardButton, 0, programsLen)

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

	nextParams := bot_types.NewEmptyParams()
	nextParams.UserId = clientId

	previousParams := bot_types.NewEmptyParams()
	previousParams.UserId = clientId

	programKb = append(programKb, GetPaginationButtons(
		programsLen,
		totalProgramCount,
		callback_data.ClientProgramAdd,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := bot_types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(callback_data.ClientSelected, backParams)),
	}
}

func ClientProgramResultsModifyList(clientId int64, exercises []models.UserExerciseRecord, totalExerciseCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	exercisesLen := len(exercises)
	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, exercisesLen)

	for _, exercise := range exercises {
		params := bot_types.NewEmptyParams()

		params.UserId = clientId
		params.UserProgramId = exercise.UserProgramId
		params.UserExerciseRecordId = exercise.Id

		text := fmt.Sprintf("%s (%d повторень)", exercise.Name(), exercise.Reps)

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         text,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultModifySelected, params),
			},
		})
	}

	nextParams := bot_types.NewEmptyParams()
	nextParams.UserId = clientId
	nextParams.UserProgramId = exercises[0].UserProgramId

	previousParams := bot_types.NewEmptyParams()
	previousParams.UserId = clientId
	previousParams.UserProgramId = exercises[0].UserProgramId

	exerciseKb = append(exerciseKb, GetPaginationButtons(
		exercisesLen,
		totalExerciseCount,
		callback_data.ClientResultModifyList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := bot_types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(exerciseKb, GetBackButton(callback_data.ClientSelected, backParams)),
	}
}