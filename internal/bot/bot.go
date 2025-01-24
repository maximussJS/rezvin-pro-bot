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

	DefaultHandler      handlers.IDefaultHandler              `name:"DefaultHandler"`
	CommandsHandler     handlers.ICommandHandler              `name:"CommandHandler"`
	RegisterHandler     callback_queries.IRegisterHandler     `name:"RegisterHandler"`
	ProgramHandler      callback_queries.IProgramHandler      `name:"ProgramHandler"`
	ExerciseHandler     callback_queries.IExerciseHandler     `name:"ExerciseHandler"`
	PendingUsersHandler callback_queries.IPendingUsersHandler `name:"PendingUsersHandler"`
	BackHandler         callback_queries.IBackHandler         `name:"BackHandler"`
	ClientHandler       callback_queries.IClientHandler       `name:"ClientHandler"`
	UserHandler         callback_queries.IUserHandler         `name:"UserHandler"`
	MainHandler         callback_queries.IMainHandler         `name:"MainHandler"`

	UserRepository               repositories.IUserRepository               `name:"UserRepository"`
	ProgramRepository            repositories.IProgramRepository            `name:"ProgramRepository"`
	UserProgramRepository        repositories.IUserProgramRepository        `name:"UserProgramRepository"`
	UserExerciseRecordRepository repositories.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
	ExerciseRepository           repositories.IExerciseRepository           `name:"ExerciseRepository"`
}

type bot struct {
	shutdownWaitGroup *sync.WaitGroup
	shutdownContext   context.Context

	logger logger.ILogger
	config config.IConfig `name:"Config"`

	bot *tg_bot.Bot

	commandsHandler     handlers.ICommandHandler
	defaultHandler      handlers.IDefaultHandler
	registerHandler     callback_queries.IRegisterHandler
	programHandler      callback_queries.IProgramHandler
	exerciseHandler     callback_queries.IExerciseHandler
	pendingUsersHandler callback_queries.IPendingUsersHandler
	backHandler         callback_queries.IBackHandler
	clientHandler       callback_queries.IClientHandler
	userHandler         callback_queries.IUserHandler
	mainHandler         callback_queries.IMainHandler

	userRepository               repositories.IUserRepository
	programRepository            repositories.IProgramRepository
	userProgramRepository        repositories.IUserProgramRepository
	userExerciseRecordRepository repositories.IUserExerciseRecordRepository
	exerciseRepository           repositories.IExerciseRepository
}

func NewBot(deps botDependencies) *bot {
	b := &bot{
		shutdownWaitGroup:            deps.ShutdownWaitGroup,
		shutdownContext:              deps.ShutdownContext,
		logger:                       deps.Logger,
		config:                       deps.Config,
		commandsHandler:              deps.CommandsHandler,
		defaultHandler:               deps.DefaultHandler,
		programHandler:               deps.ProgramHandler,
		registerHandler:              deps.RegisterHandler,
		exerciseHandler:              deps.ExerciseHandler,
		pendingUsersHandler:          deps.PendingUsersHandler,
		backHandler:                  deps.BackHandler,
		userHandler:                  deps.UserHandler,
		mainHandler:                  deps.MainHandler,
		clientHandler:                deps.ClientHandler,
		userRepository:               deps.UserRepository,
		programRepository:            deps.ProgramRepository,
		userProgramRepository:        deps.UserProgramRepository,
		userExerciseRecordRepository: deps.UserExerciseRecordRepository,
		exerciseRepository:           deps.ExerciseRepository,
	}

	opts := []tg_bot.Option{
		tg_bot.WithDefaultHandler(b.defaultHandler.Handle),
		tg_bot.WithMiddlewares(b.panicRecoveryMiddleware, b.chatIdMiddleware),
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

func (bot *bot) registerCommand(command string, handler tg_bot.HandlerFunc, middlewares []tg_bot.Middleware) {
	bot.bot.RegisterHandler(
		tg_bot.HandlerTypeMessageText,
		command,
		tg_bot.MatchTypeExact,
		handler,
		middlewares...,
	)
}

func (bot *bot) registerCallbackQueryByPrefix(prefix string, handler tg_bot.HandlerFunc, middlewares []tg_bot.Middleware) {
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

	bot.registerCommand(constants.CommandStart, bot.commandsHandler.Start, bot.commandMiddlewares())

	bot.registerCallbackQueryByPrefix(callback_data.MainPrefix, bot.mainHandler.Handle, bot.mainMiddlewares())

	bot.registerCallbackQueryByPrefix(callback_data.RegisterPrefix, bot.registerHandler.Handle, bot.userMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.UserPrefix, bot.userHandler.Handle, bot.userMiddlewares())

	bot.registerCallbackQueryByPrefix(callback_data.ProgramPrefix, bot.programHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.ExercisePrefix, bot.exerciseHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.PendingUsersPrefix, bot.pendingUsersHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.BackPrefix, bot.backHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.ClientPrefix, bot.clientHandler.Handle, bot.adminMiddlewares())
}
