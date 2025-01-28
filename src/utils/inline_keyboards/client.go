package inline_keyboards

import (
	"fmt"
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants/callback_data"
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
		types.NewEmptyParams(),
		types.NewEmptyParams(),
	))

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(clientKb, GetBackButton(callback_data.MainBackToMain, types.NewEmptyParams())),
	}
}

func ClientSelectedMenu(clientId int64) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

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
	params := types.NewEmptyParams()

	params.UserId = clientId
	params.UserProgramId = program.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultList, params)},
			},
			{
				{Text: "✍️ Внести результати", CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultModifyExercisesList, params)},
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
		params := types.NewEmptyParams()

		params.UserId = clientId
		params.UserProgramId = program.Id

		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name(),
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramSelected, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()
	nextParams.UserId = clientId

	previousParams := types.NewEmptyParams()
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

	backParams := types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(callback_data.ClientSelected, backParams)),
	}
}

func ProgramForClientList(clientId int64, programs []models.Program, totalProgramCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	programsLen := len(programs)
	programKb := make([][]tg_models.InlineKeyboardButton, 0, programsLen)

	for _, program := range programs {
		params := types.NewEmptyParams()

		params.UserId = clientId
		params.ProgramId = program.Id
		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientProgramAssign, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()
	nextParams.UserId = clientId

	previousParams := types.NewEmptyParams()
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

	backParams := types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(callback_data.ClientSelected, backParams)),
	}
}

func ClientProgramResultsModifyExerciseList(clientId int64, userProgramId uint, exercises []models.Exercise, totalExerciseCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	exercisesLen := len(exercises)
	exerciseKb := make([][]tg_models.InlineKeyboardButton, 0, exercisesLen)

	for _, exercise := range exercises {
		params := types.NewEmptyParams()

		params.UserId = clientId
		params.UserProgramId = userProgramId
		params.ExerciseId = exercise.Id

		exerciseKb = append(exerciseKb, []tg_models.InlineKeyboardButton{
			{
				Text:         exercise.Name,
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultModifyExerciseSelected, params),
			},
		})
	}

	nextParams := types.NewEmptyParams()
	nextParams.UserId = clientId
	nextParams.UserProgramId = userProgramId

	previousParams := types.NewEmptyParams()
	previousParams.UserId = clientId
	previousParams.UserProgramId = userProgramId

	exerciseKb = append(exerciseKb, GetPaginationButtons(
		exercisesLen,
		totalExerciseCount,
		callback_data.ClientResultModifyExercisesList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(exerciseKb, GetBackButton(callback_data.ClientSelected, backParams)),
	}
}

func ClientSelectedOk(clientId int64) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()
	params.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(callback_data.ClientSelected, params),
		},
	}
}

func ClientProgramSelectedOk(clientId int64, userProgramId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()
	params.UserId = clientId
	params.UserProgramId = userProgramId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(callback_data.ClientProgramSelected, params),
		},
	}
}

func ClientProgramResultModifyExerciseSelectedOk(clientId int64, records []models.UserExerciseRecord) *tg_models.InlineKeyboardMarkup {
	recordsLen := len(records)
	recordsKb := make([][]tg_models.InlineKeyboardButton, 0, recordsLen)

	for _, record := range records {
		params := types.NewEmptyParams()

		params.UserId = clientId
		params.UserProgramId = record.UserProgramId
		params.ExerciseId = record.ExerciseId
		params.UserExerciseRecordId = record.Id

		recordsKb = append(recordsKb, []tg_models.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d повторень", record.Reps),
				CallbackData: bot_utils.AddParamsToQueryString(callback_data.ClientResultModifyExerciseRepsModify, params),
			},
		})
	}

	backParams := types.NewEmptyParams()
	backParams.UserId = clientId
	backParams.UserProgramId = records[0].UserProgramId
	backParams.ExerciseId = records[0].ExerciseId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(recordsKb, GetBackButton(callback_data.ClientResultModifyExercisesList, backParams)),
	}
}
