package services

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/internal/logger"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/repositories"
	"rezvin-pro-bot/src/utils"
)

type ISenderService interface {
	AnswerCallbackQuery(ctx context.Context, b *tg_bot.Bot, update *tg_models.Update) bool
	Send(ctx context.Context, b *tg_bot.Bot, chatId int64, message string) int
	SendSafe(ctx context.Context, b *tg_bot.Bot, chatId int64, message string) int
	SendWithKb(
		ctx context.Context,
		b *tg_bot.Bot,
		chatId int64,
		message string,
		kb *tg_models.InlineKeyboardMarkup,
	) int
	SendSafeWithKb(
		ctx context.Context,
		b *tg_bot.Bot,
		chatId int64,
		message string,
		kb *tg_models.InlineKeyboardMarkup,
	) int
	Delete(ctx context.Context, b *tg_bot.Bot, chatId int64, messageId int)
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
) int {
	return s.send(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        message,
		ReplyMarkup: kb,
		ParseMode:   tg_models.ParseModeMarkdown,
	}, false)
}

func (s *senderService) SendSafeWithKb(
	ctx context.Context,
	b *tg_bot.Bot,
	chatId int64,
	message string,
	kb *tg_models.InlineKeyboardMarkup,
) int {
	return s.send(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      chatId,
		Text:        message,
		ReplyMarkup: kb,
		ParseMode:   tg_models.ParseModeMarkdown,
	}, true)
}

func (s *senderService) Send(ctx context.Context, b *tg_bot.Bot, chatId int64, message string) int {
	return s.send(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      message,
		ParseMode: tg_models.ParseModeMarkdown,
	}, false)
}

func (s *senderService) SendSafe(ctx context.Context, b *tg_bot.Bot, chatId int64, message string) int {
	return s.send(ctx, b, &tg_bot.SendMessageParams{
		ChatID:    chatId,
		Text:      message,
		ParseMode: tg_models.ParseModeMarkdown,
	}, true)
}

func (s *senderService) send(ctx context.Context, b *tg_bot.Bot, params *tg_bot.SendMessageParams, safe bool) int {
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

	return msg.ID
}

func (s *senderService) Delete(ctx context.Context, b *tg_bot.Bot, chatId int64, messageId int) {
	lastMsg := s.lastMessageRepository.GetByChatId(ctx, chatId)

	if lastMsg == nil && lastMsg.MessageId == messageId {
		s.lastMessageRepository.DeleteByChatId(ctx, chatId)
	}

	ok, err := b.DeleteMessage(ctx, &tg_bot.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: messageId,
	})

	if err != nil || !ok {
		s.logger.Error(err.Error())
	}
}
