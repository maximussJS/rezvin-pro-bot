package callback_queries

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/constants"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/services"
	bot_utils "rezvin-pro-bot/utils/bot"
	"strings"
)

type IClientHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type clientHandlerDependencies struct {
	dig.In

	Logger                       logger.ILogger                             `name:"Logger"`
	TextService                  services.ITextService                      `name:"TextService"`
	InlineKeyboardService        services.IInlineKeyboardService            `name:"InlineKeyboardService"`
	UserRepository               repositories.IUserRepository               `name:"UserRepository"`
	ProgramRepository            repositories.IProgramRepository            `name:"ProgramRepository"`
	UserProgramRepository        repositories.IUserProgramRepository        `name:"UserProgramRepository"`
	UserExerciseRecordRepository repositories.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
}

type clientHandler struct {
	logger                       logger.ILogger
	textService                  services.ITextService
	inlineKeyboardService        services.IInlineKeyboardService
	userRepository               repositories.IUserRepository
	programRepository            repositories.IProgramRepository
	userProgramRepository        repositories.IUserProgramRepository
	userExerciseRecordRepository repositories.IUserExerciseRecordRepository
}

func NewClientHandler(deps clientHandlerDependencies) *clientHandler {
	return &clientHandler{
		logger:                       deps.Logger,
		textService:                  deps.TextService,
		inlineKeyboardService:        deps.InlineKeyboardService,
		userRepository:               deps.UserRepository,
		programRepository:            deps.ProgramRepository,
		userProgramRepository:        deps.UserProgramRepository,
		userExerciseRecordRepository: deps.UserExerciseRecordRepository,
	}
}

func (h *clientHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	callBackQueryData := update.CallbackQuery.Data

	if strings.HasPrefix(callBackQueryData, callback_data.ClientSelected) {
		h.selected(ctx, b, update)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientProgramList) {
		h.programList(ctx, b, update)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientProgramAdd) {
		h.programAdd(ctx, b, update)
		return
	}

	if strings.HasPrefix(callBackQueryData, callback_data.ClientProgramAssign) {
		h.programAssign(ctx, b, update)
		return
	}

	switch callBackQueryData {
	case callback_data.ClientList:
		h.list(ctx, b, update)
	}
}

func (h *clientHandler) list(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)

	clients := h.userRepository.GetClients(ctx, 5, 0)

	if len(clients) == 0 {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.NoClientsMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectClientMessage(),
		ReplyMarkup: h.inlineKeyboardService.ClientList(clients),
		ParseMode:   tg_models.ParseModeMarkdown,
	})
}

func (h *clientHandler) selected(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	userId := bot_utils.GetSelectedUserId(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		h.logger.Error(fmt.Sprintf("User not found: %d", userId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectClientOptionMessage(user.GetReadableName()),
		ParseMode:   tg_models.ParseModeMarkdown,
		ReplyMarkup: h.inlineKeyboardService.ClientSelectedMenu(userId),
	})
}

func (h *clientHandler) programList(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	userId := bot_utils.GetSelectedUserId(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		h.logger.Error(fmt.Sprintf("User not found: %d", userId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	programs := h.userProgramRepository.GetByUserId(ctx, userId, 5, 0)

	if len(programs) == 0 {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.NoClientProgramsMessage(user.GetReadableName()),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectClientProgramMessage(user.GetReadableName()),
		ParseMode:   tg_models.ParseModeMarkdown,
		ReplyMarkup: h.inlineKeyboardService.ClientProgramList(userId, programs),
	})
}

func (h *clientHandler) programAdd(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	userId := bot_utils.GetSelectedUserId(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		h.logger.Error(fmt.Sprintf("User not found: %d", userId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	programs := h.programRepository.GetNotAssignedToUser(ctx, userId, 5, 0)

	if len(programs) == 0 {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.NoProgramsForClientMessage(user.GetReadableName()),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectClientProgramMessage(user.GetReadableName()),
		ParseMode:   tg_models.ParseModeMarkdown,
		ReplyMarkup: h.inlineKeyboardService.ProgramForClientList(userId, programs),
	})
}

func (h *clientHandler) programAssign(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	userId := bot_utils.GetSelectedUserId(update)
	programId := bot_utils.GetClientProgramId(update)

	user := h.userRepository.GetById(ctx, userId)

	if user == nil {
		h.logger.Error(fmt.Sprintf("User not found: %d", userId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	userProgram := h.userProgramRepository.GetByUserIdAndProgramId(ctx, userId, programId)

	if userProgram != nil {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ClientProgramAlreadyAssignedMessage(user.GetReadableName(), userProgram.Program.Name),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		h.logger.Error(fmt.Sprintf("Program with id %d not found", programId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	h.userProgramRepository.Create(ctx, models.UserProgram{
		UserId:    userId,
		ProgramId: programId,
	})

	records := make([]models.UserExerciseRecord, 0, len(program.Exercises))

	for _, exercise := range program.Exercises {
		for _, rep := range constants.RepsList {
			records = append(records, models.UserExerciseRecord{
				UserProgramId: programId,
				ExerciseId:    exercise.Id,
				Weight:        0,
				Reps:          uint(rep),
			})
		}
	}

	h.userExerciseRecordRepository.CreateMany(ctx, records)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.ClientProgramAssignedMessage(user.GetReadableName(), program.Name),
		ParseMode: tg_models.ParseModeMarkdown,
	})
}
