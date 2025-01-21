package callback_queries

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/services"
	bot_utils "rezvin-pro-bot/utils/bot"
	"strings"
)

type IExerciseHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update)
}

type exerciseHandlerDependencies struct {
	dig.In

	Logger                logger.ILogger                  `name:"Logger"`
	TextService           services.ITextService           `name:"TextService"`
	ConversationService   services.IConversationService   `name:"ConversationService"`
	InlineKeyboardService services.IInlineKeyboardService `name:"InlineKeyboardService"`

	ProgramRepository  repositories.IProgramRepository  `name:"ProgramRepository"`
	ExerciseRepository repositories.IExerciseRepository `name:"ExerciseRepository"`
}

type exerciseHandler struct {
	logger                logger.ILogger
	textService           services.ITextService
	conversationService   services.IConversationService
	inlineKeyboardService services.IInlineKeyboardService
	programRepository     repositories.IProgramRepository
	exerciseRepository    repositories.IExerciseRepository
}

func NewExerciseHandler(deps exerciseHandlerDependencies) *exerciseHandler {
	return &exerciseHandler{
		logger:                deps.Logger,
		textService:           deps.TextService,
		inlineKeyboardService: deps.InlineKeyboardService,
		conversationService:   deps.ConversationService,
		programRepository:     deps.ProgramRepository,
		exerciseRepository:    deps.ExerciseRepository,
	}
}

func (h *exerciseHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	answerResult := bot_utils.MustAnswerCallbackQuery(ctx, b, update)

	if !answerResult {
		h.logger.Error(fmt.Sprintf("Failed to answer callback query: %s", update.CallbackQuery.ID))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   h.textService.ErrorMessage(),
		})
		return
	}

	callbackDataQuery := update.CallbackQuery.Data

	if strings.HasPrefix(callbackDataQuery, callback_data.ExerciseAdd) {
		h.add(ctx, b, update)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ExerciseList) {
		h.list(ctx, b, update)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ExerciseDeleteItem) {
		h.deleteItem(ctx, b, update)
		return
	}

	if strings.HasPrefix(callbackDataQuery, callback_data.ExerciseDelete) {
		h.delete(ctx, b, update)
		return
	}
}

func (h *exerciseHandler) add(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		h.logger.Error(fmt.Sprintf("Program not found: %d", programId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	conversation := h.conversationService.CreateConversation(chatId)
	defer h.conversationService.DeleteConversation(chatId)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.EnterExerciseNameMessage(),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	exerciseName := conversation.WaitAnswer()

	existingExercise := h.exerciseRepository.GetByNameAndProgramId(ctx, exerciseName, programId)

	if existingExercise != nil {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ExerciseNameAlreadyExistsMessage(exerciseName),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	h.exerciseRepository.Create(ctx, models.Exercise{
		Name:      exerciseName,
		ProgramId: programId,
	})

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.ExerciseSuccessfullyAddedMessage(exerciseName, program.Name),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	h.backToSelectedProgram(ctx, b, update, program)
}

func (h *exerciseHandler) list(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		h.logger.Error(fmt.Sprintf("Program not found: %d", programId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	if len(program.Exercises) == 0 {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.NoExercisesMessage(program.Name),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.ExercisesMessage(program.Name, program.Exercises),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	h.backToSelectedProgram(ctx, b, update, program)
}

func (h *exerciseHandler) delete(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		h.logger.Error(fmt.Sprintf("Program not found: %d", programId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	if len(program.Exercises) == 0 {
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.NoExercisesMessage(program.Name),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.ExerciseDeleteMessage(program.Name),
		ParseMode:   tg_models.ParseModeMarkdown,
		ReplyMarkup: h.inlineKeyboardService.ProgramExerciseDeleteList(program.Id, program.Exercises),
	})
}

func (h *exerciseHandler) deleteItem(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) {
	chatId := bot_utils.GetChatID(update)
	programId := bot_utils.GetProgramId(update)
	exerciseId := bot_utils.GetExerciseId(update)

	program := h.programRepository.GetById(ctx, programId)

	if program == nil {
		h.logger.Error(fmt.Sprintf("Program not found: %d", programId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	exercise := h.exerciseRepository.GetByIdAndProgramId(ctx, exerciseId, programId)

	if exercise == nil {
		h.logger.Error(fmt.Sprintf("Exercise not found: %d", exerciseId))
		bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    chatId,
			Text:      h.textService.ErrorMessage(),
			ParseMode: tg_models.ParseModeMarkdown,
		})
		return
	}

	h.exerciseRepository.DeleteById(ctx, exerciseId)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      h.textService.ExerciseSuccessfullyDeletedMessage(exercise.Name, program.Name),
		ParseMode: tg_models.ParseModeMarkdown,
	})

	h.backToSelectedProgram(ctx, b, update, program)
}

func (h *exerciseHandler) backToSelectedProgram(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update, program *models.Program) {
	chatId := bot_utils.GetChatID(update)

	bot_utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        h.textService.SelectProgramOptionMessage(program.Name),
		ParseMode:   tg_models.ParseModeMarkdown,
		ReplyMarkup: h.inlineKeyboardService.ProgramSelectedMenu(program.Id),
	})
}
