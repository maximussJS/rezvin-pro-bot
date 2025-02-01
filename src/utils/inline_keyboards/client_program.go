package inline_keyboards

import (
	tg_models "github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/types"
	bot_utils "rezvin-pro-bot/src/utils/bot"
)

func ClientProgramMenu(clientId int64, program models.UserProgram) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()

	params.UserId = clientId
	params.UserProgramId = program.Id

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			{
				{Text: "🚀 Переглянути результати", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientResultList, params)},
			},
			{
				{Text: "✍️ Внести результати", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientResultExercisesList, params)},
			},
			{
				{Text: "➖ Видалити програму", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientProgramDelete, params)},
			},
			{
				{Text: "🔙 Назад", CallbackData: bot_utils.AddParamsToQueryString(constants.ClientProgramList, params)},
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
				CallbackData: bot_utils.AddParamsToQueryString(constants.ClientProgramSelected, params),
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
		constants.ClientProgramList,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(constants.ClientSelected, backParams)),
	}
}

func ClientProgramAssignList(clientId int64, programs []models.Program, totalProgramCount int64, limit, offset int) *tg_models.InlineKeyboardMarkup {
	programsLen := len(programs)
	programKb := make([][]tg_models.InlineKeyboardButton, 0, programsLen)

	for _, program := range programs {
		params := types.NewEmptyParams()

		params.UserId = clientId
		params.ProgramId = program.Id
		programKb = append(programKb, []tg_models.InlineKeyboardButton{
			{
				Text:         program.Name,
				CallbackData: bot_utils.AddParamsToQueryString(constants.ClientProgramAssign, params),
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
		constants.ClientProgramAdd,
		limit,
		offset,
		nextParams,
		previousParams,
	))

	backParams := types.NewEmptyParams()

	backParams.UserId = clientId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: append(programKb, GetBackButton(constants.ClientSelected, backParams)),
	}
}

func ClientProgramSelectedOk(clientId int64, userProgramId uint) *tg_models.InlineKeyboardMarkup {
	params := types.NewEmptyParams()
	params.UserId = clientId
	params.UserProgramId = userProgramId

	return &tg_models.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg_models.InlineKeyboardButton{
			GetOkButton(constants.ClientProgramSelected, params),
		},
	}
}
