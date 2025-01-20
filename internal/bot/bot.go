package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"go.uber.org/dig"
	"rezvin-pro-bot/config"
	"rezvin-pro-bot/constants"
	"rezvin-pro-bot/constants/callback_data"
	"rezvin-pro-bot/handlers"
	"rezvin-pro-bot/handlers/callback_queries"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/services"
	"rezvin-pro-bot/utils"
	"sync"
)

type IBot interface {
	Start()
}

type botDependencies struct {
	dig.In

	ShutdownWaitGroup *sync.WaitGroup `name:"ShutdownWaitGroup"`
	ShutdownContext   context.Context `name:"ShutdownContext"`
	Logger            logger.ILogger  `name:"Logger"`
	Config            config.IConfig  `name:"Config"`

	DefaultHandler  handlers.IDefaultHandler         `name:"DefaultHandler"`
	CommandsHandler handlers.ICommandHandler         `name:"CommandHandler"`
	UserHandler     callback_queries.IUserHandler    `name:"UserHandler"`
	AdminHandler    callback_queries.IAdminHandler   `name:"AdminHandler"`
	ProgramHandler  callback_queries.IProgramHandler `name:"ProgramHandler"`

	TextService    services.ITextService        `name:"TextService"`
	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type bot struct {
	shutdownWaitGroup *sync.WaitGroup
	shutdownContext   context.Context

	logger logger.ILogger
	config config.IConfig `name:"Config"`

	bot *tg_bot.Bot

	commandsHandler handlers.ICommandHandler
	defaultHandler  handlers.IDefaultHandler
	userHandler     callback_queries.IUserHandler
	adminHandler    callback_queries.IAdminHandler
	programHandler  callback_queries.IProgramHandler

	textService    services.ITextService
	userRepository repositories.IUserRepository
}

func NewBot(deps botDependencies) *bot {
	b := &bot{
		shutdownWaitGroup: deps.ShutdownWaitGroup,
		shutdownContext:   deps.ShutdownContext,
		logger:            deps.Logger,
		config:            deps.Config,
		commandsHandler:   deps.CommandsHandler,
		defaultHandler:    deps.DefaultHandler,
		adminHandler:      deps.AdminHandler,
		programHandler:    deps.ProgramHandler,
		userHandler:       deps.UserHandler,
		textService:       deps.TextService,
		userRepository:    deps.UserRepository,
	}

	opts := []tg_bot.Option{
		tg_bot.WithDefaultHandler(b.defaultHandler.Handle),
		tg_bot.WithMiddlewares(b.panicRecoveryMiddleware),
	}

	tgBot, err := tg_bot.New(b.config.BotToken(), opts...)

	utils.PanicIfError(err)

	b.bot = tgBot

	b.registerHandlers()

	return b
}

func (bot *bot) Start() {
	defer bot.shutdownWaitGroup.Done()

	bot.shutdownWaitGroup.Add(1)

	bot.logger.Log("Bot started")

	bot.bot.Start(bot.shutdownContext)

	select {
	case <-bot.shutdownContext.Done():
		bot.logger.Log("Shutting down Bot gracefully...")
	}

	bot.logger.Log("Bot stopped")
}

func (bot *bot) registerCommand(command string, handler tg_bot.HandlerFunc, middlewares ...tg_bot.Middleware) {
	bot.bot.RegisterHandler(
		tg_bot.HandlerTypeMessageText,
		command,
		tg_bot.MatchTypeExact,
		handler,
		middlewares...,
	)
}

func (bot *bot) registerCallbackQueryByPrefix(prefix string, handler tg_bot.HandlerFunc, middlewares ...tg_bot.Middleware) {
	bot.bot.RegisterHandler(
		tg_bot.HandlerTypeCallbackQueryData,
		prefix,
		tg_bot.MatchTypePrefix,
		handler,
		middlewares...,
	)
}

func (bot *bot) registerHandlers() {
	if bot.bot == nil {
		panic("cannot register handlers without bot instance")
	}

	bot.registerCommand(constants.CommandStart, bot.commandsHandler.Start, bot.timeoutMiddleware)

	bot.registerCallbackQueryByPrefix(callback_data.UserPrefix, bot.userHandler.Handle, bot.timeoutMiddleware)
	bot.registerCallbackQueryByPrefix(callback_data.AdminPrefix, bot.adminHandler.Handle, bot.isAdminMiddleware)
	bot.registerCallbackQueryByPrefix(callback_data.ProgramPrefix, bot.programHandler.Handle, bot.isAdminMiddleware)
}
