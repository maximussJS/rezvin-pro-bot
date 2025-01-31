package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"go.uber.org/dig"
	"net/http"
	"rezvin-pro-bot/src/config"
	"rezvin-pro-bot/src/handlers"
	"rezvin-pro-bot/src/handlers/callback_queries"
	"rezvin-pro-bot/src/internal/logger"
	"rezvin-pro-bot/src/repositories"
	"rezvin-pro-bot/src/services"
	"rezvin-pro-bot/src/utils"
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

	SenderService       services.ISenderService       `name:"SenderService"`
	LockService         services.ILockService         `name:"LockService"`
	ConversationService services.IConversationService `name:"ConversationService"`

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

	bot    *tg_bot.Bot
	server http.Server

	senderService       services.ISenderService
	lockService         services.ILockService
	conversationService services.IConversationService

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
		shutdownWaitGroup: deps.ShutdownWaitGroup,
		shutdownContext:   deps.ShutdownContext,
		logger:            deps.Logger,
		config:            deps.Config,

		senderService:       deps.SenderService,
		lockService:         deps.LockService,
		conversationService: deps.ConversationService,

		commandsHandler:     deps.CommandsHandler,
		defaultHandler:      deps.DefaultHandler,
		programHandler:      deps.ProgramHandler,
		registerHandler:     deps.RegisterHandler,
		exerciseHandler:     deps.ExerciseHandler,
		pendingUsersHandler: deps.PendingUsersHandler,
		backHandler:         deps.BackHandler,
		userHandler:         deps.UserHandler,
		mainHandler:         deps.MainHandler,
		clientHandler:       deps.ClientHandler,

		userRepository:               deps.UserRepository,
		programRepository:            deps.ProgramRepository,
		userProgramRepository:        deps.UserProgramRepository,
		userExerciseRecordRepository: deps.UserExerciseRecordRepository,
		exerciseRepository:           deps.ExerciseRepository,
	}

	opts := []tg_bot.Option{
		tg_bot.WithWebhookSecretToken(b.config.WebhookSecretToken()),
		tg_bot.WithDefaultHandler(b.defaultHandler.Handle),
		tg_bot.WithMiddlewares(b.defaultMiddlewares()...),
	}

	tgBot, err := tg_bot.New(b.config.BotToken(), opts...)

	utils.PanicIfError(err)

	b.bot = tgBot

	b.registerHandlers()

	return b
}
