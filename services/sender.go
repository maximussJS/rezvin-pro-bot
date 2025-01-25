package services

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/utils"
)

type ISenderService interface {
	AnswerCallbackQuery(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) bool
	Send(ctx context.Context, b *tg_bot.Bot, chatId int64, message string)
	SendSafe(ctx context.Context, b *tg_bot.Bot, chatId int64, message string)
	SendWithKb(
		ctx context.Context,
		b *tg_bot.Bot,
		chatId int64,
		message string,
		kb *tg_models.InlineKeyboardMarkup,
	)
	SendSafeWithKb(
		ctx context.Context,
		b *tg_bot.Bot,
		chatId int64,
		message string,
		kb *tg_models.InlineKeyboardMarkup,
	)
}

type senderServiceDependencies struct {
	dig.In

	Logger                logger.ILogger                          `name:"Logger"`
	LastMessageRepository repositories.ILastUserMessageRepository `name:"LastUserMessageRepository"`
}

type senderService struct {
	logger                logger.ILogger
	lastMessageRepository repositories.ILastUserMessageRepository
}

func NewSenderService(deps senderServiceDependencies) *senderService {
	return &senderService{
		logger:                deps.Logger,
		lastMessageRepository: deps.LastMessageRepository,
	}
}

func (s *senderService) AnswerCallbackQuery(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) bool {
	result, err := b.AnswerCallbackQuery(ctx, &tg_bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	utils.PanicIfError(err)

	return result
}

func (s *senderService) SendWithKb(
	ctx context.Context,
	b *tg_bot.Bot,
	chatId int64,
	message string,
	kb *tg_models.InlineKeyboardMarkup,
) {
	s.send(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        message,
		ReplyMarkup: kb,
		ParseMode:   tg_models.ParseModeMarkdown,
	}, true)
}

func (s *senderService) SendSafeWithKb(
	ctx context.Context,
	b *tg_bot.Bot,
	chatId int64,
	message string,
	kb *tg_models.InlineKeyboardMarkup,
) {
	s.send(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        message,
		ReplyMarkup: kb,
		ParseMode:   tg_models.ParseModeMarkdown,
	}, false)
}

func (s *senderService) Send(ctx context.Context, b *tg_bot.Bot, chatId int64, message string) {
	s.send(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      message,
		ParseMode: tg_models.ParseModeMarkdown,
	}, true)
}

func (s *senderService) SendSafe(ctx context.Context, b *tg_bot.Bot, chatId int64, message string) {
	s.send(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      message,
		ParseMode: tg_models.ParseModeMarkdown,
	}, false)
}

func (s *senderService) send(ctx context.Context, b *tg_bot.Bot, params *tg_bot.SendMessageParams, safe bool) {
	chatId := params.ChatID.(int64)

	lastMsg := s.lastMessageRepository.GetByChatId(ctx, chatId)

	if lastMsg != nil && !safe {
		ok, err := b.DeleteMessage(ctx, &tg_bot.DeleteMessageParams{
			ChatID:    chatId,
			MessageID: lastMsg.MessageId,
		})

		if err != nil || !ok {
			s.logger.Error(err.Error())
			s.lastMessageRepository.DeleteByChatId(ctx, chatId)
		}
	}

	msg, err := b.SendMessage(ctx, params)

	utils.PanicIfError(err)

	if !safe {
		if lastMsg == nil {
			s.lastMessageRepository.Create(ctx, models.LastUserMessage{
				ChatId:    chatId,
				MessageId: msg.ID,
			})
		} else {
			s.lastMessageRepository.UpdateByChatId(ctx, chatId, models.LastUserMessage{
				MessageId: msg.ID,
			})
		}
	}
}
