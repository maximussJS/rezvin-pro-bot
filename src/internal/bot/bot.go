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
)

type IBot interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context) error
}

type botDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`
	Config config.IConfig `name:"Config"`

	SenderService       services.ISenderService       `name:"SenderService"`
	LockService         services.ILockService         `name:"LockService"`
	ConversationService services.IConversationService `name:"ConversationService"`

	DefaultHandler       handlers.IDefaultHandler               `name:"DefaultHandler"`
	CommandsHandler      handlers.ICommandHandler               `name:"CommandHandler"`
	RegisterHandler      callback_queries.IRegisterHandler      `name:"RegisterHandler"`
	ProgramHandler       callback_queries.IProgramHandler       `name:"ProgramHandler"`
	ExerciseHandler      callback_queries.IExerciseHandler      `name:"ExerciseHandler"`
	MeasureHandler       callback_queries.IMeasureHandler       `name:"MeasureHandler"`
	PendingUsersHandler  callback_queries.IPendingUsersHandler  `name:"PendingUsersHandler"`
	BackHandler          callback_queries.IBackHandler          `name:"BackHandler"`
	ClientHandler        callback_queries.IClientHandler        `name:"ClientHandler"`
	ClientProgramHandler callback_queries.IClientProgramHandler `name:"ClientProgramHandler"`
	ClientResultHandler  callback_queries.IClientResultHandler  `name:"ClientResultHandler"`
	ClientMeasureHandler callback_queries.IClientMeasureHandler `name:"ClientMeasureHandler"`
	UserResultHandler    callback_queries.IUserResultHandler    `name:"UserResultHandler"`
	UserProgramHandler   callback_queries.IUserProgramHandler   `name:"UserProgramHandler"`
	MainHandler          callback_queries.IMainHandler          `name:"MainHandler"`

	UserRepository        repositories.IUserRepository        `name:"UserRepository"`
	ProgramRepository     repositories.IProgramRepository     `name:"ProgramRepository"`
	UserProgramRepository repositories.IUserProgramRepository `name:"UserProgramRepository"`
	UserMeasureRepository repositories.IUserMeasureRepository `name:"UserMeasureRepository"`
	UserResultRepository  repositories.IUserResultRepository  `name:"UserResultRepository"`
	ExerciseRepository    repositories.IExerciseRepository    `name:"ExerciseRepository"`
	MeasureRepository     repositories.IMeasureRepository     `name:"MeasureRepository"`
}

type bot struct {
	logger logger.ILogger
	config config.IConfig `name:"Config"`

	bot    *tg_bot.Bot
	server http.Server

	senderService       services.ISenderService
	lockService         services.ILockService
	conversationService services.IConversationService

	commandsHandler      handlers.ICommandHandler
	defaultHandler       handlers.IDefaultHandler
	registerHandler      callback_queries.IRegisterHandler
	programHandler       callback_queries.IProgramHandler
	exerciseHandler      callback_queries.IExerciseHandler
	measureHandler       callback_queries.IMeasureHandler
	pendingUsersHandler  callback_queries.IPendingUsersHandler
	backHandler          callback_queries.IBackHandler
	clientHandler        callback_queries.IClientHandler
	clientProgramHandler callback_queries.IClientProgramHandler
	clientResultHandler  callback_queries.IClientResultHandler
	clientMeasureHandler callback_queries.IClientMeasureHandler
	userResultHandler    callback_queries.IUserResultHandler
	userProgramHandler   callback_queries.IUserProgramHandler
	mainHandler          callback_queries.IMainHandler

	userRepository        repositories.IUserRepository
	programRepository     repositories.IProgramRepository
	userProgramRepository repositories.IUserProgramRepository
	userResultRepository  repositories.IUserResultRepository
	userMeasureRepository repositories.IUserMeasureRepository
	exerciseRepository    repositories.IExerciseRepository
	measureRepository     repositories.IMeasureRepository
}

func NewBot(deps botDependencies) *bot {
	b := &bot{
		logger: deps.Logger,
		config: deps.Config,

		senderService:       deps.SenderService,
		lockService:         deps.LockService,
		conversationService: deps.ConversationService,

		commandsHandler:      deps.CommandsHandler,
		defaultHandler:       deps.DefaultHandler,
		programHandler:       deps.ProgramHandler,
		registerHandler:      deps.RegisterHandler,
		exerciseHandler:      deps.ExerciseHandler,
		measureHandler:       deps.MeasureHandler,
		pendingUsersHandler:  deps.PendingUsersHandler,
		backHandler:          deps.BackHandler,
		userResultHandler:    deps.UserResultHandler,
		userProgramHandler:   deps.UserProgramHandler,
		mainHandler:          deps.MainHandler,
		clientHandler:        deps.ClientHandler,
		clientProgramHandler: deps.ClientProgramHandler,
		clientResultHandler:  deps.ClientResultHandler,
		clientMeasureHandler: deps.ClientMeasureHandler,

		userRepository:        deps.UserRepository,
		programRepository:     deps.ProgramRepository,
		userProgramRepository: deps.UserProgramRepository,
		userResultRepository:  deps.UserResultRepository,
		userMeasureRepository: deps.UserMeasureRepository,
		exerciseRepository:    deps.ExerciseRepository,
		measureRepository:     deps.MeasureRepository,
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
