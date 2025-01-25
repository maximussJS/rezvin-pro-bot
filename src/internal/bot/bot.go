package bot

import (
	"context"
	"crypto/tls"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"go.uber.org/dig"
	"net/http"
	"rezvin-pro-bot/src/config"
	constants2 "rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/constants/callback_data"
	handlers2 "rezvin-pro-bot/src/handlers"
	callback_queries2 "rezvin-pro-bot/src/handlers/callback_queries"
	"rezvin-pro-bot/src/internal/logger"
	repositories2 "rezvin-pro-bot/src/repositories"
	"rezvin-pro-bot/src/services"
	"rezvin-pro-bot/src/utils"
	"sync"
	"time"
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

	SenderService services.ISenderService `name:"SenderService"`

	DefaultHandler      handlers2.IDefaultHandler              `name:"DefaultHandler"`
	CommandsHandler     handlers2.ICommandHandler              `name:"CommandHandler"`
	RegisterHandler     callback_queries2.IRegisterHandler     `name:"RegisterHandler"`
	ProgramHandler      callback_queries2.IProgramHandler      `name:"ProgramHandler"`
	ExerciseHandler     callback_queries2.IExerciseHandler     `name:"ExerciseHandler"`
	PendingUsersHandler callback_queries2.IPendingUsersHandler `name:"PendingUsersHandler"`
	BackHandler         callback_queries2.IBackHandler         `name:"BackHandler"`
	ClientHandler       callback_queries2.IClientHandler       `name:"ClientHandler"`
	UserHandler         callback_queries2.IUserHandler         `name:"UserHandler"`
	MainHandler         callback_queries2.IMainHandler         `name:"MainHandler"`

	UserRepository               repositories2.IUserRepository               `name:"UserRepository"`
	ProgramRepository            repositories2.IProgramRepository            `name:"ProgramRepository"`
	UserProgramRepository        repositories2.IUserProgramRepository        `name:"UserProgramRepository"`
	UserExerciseRecordRepository repositories2.IUserExerciseRecordRepository `name:"UserExerciseRecordRepository"`
	ExerciseRepository           repositories2.IExerciseRepository           `name:"ExerciseRepository"`
}

type bot struct {
	shutdownWaitGroup *sync.WaitGroup
	shutdownContext   context.Context

	logger logger.ILogger
	config config.IConfig `name:"Config"`

	bot    *tg_bot.Bot
	server http.Server

	senderService services.ISenderService

	commandsHandler     handlers2.ICommandHandler
	defaultHandler      handlers2.IDefaultHandler
	registerHandler     callback_queries2.IRegisterHandler
	programHandler      callback_queries2.IProgramHandler
	exerciseHandler     callback_queries2.IExerciseHandler
	pendingUsersHandler callback_queries2.IPendingUsersHandler
	backHandler         callback_queries2.IBackHandler
	clientHandler       callback_queries2.IClientHandler
	userHandler         callback_queries2.IUserHandler
	mainHandler         callback_queries2.IMainHandler

	userRepository               repositories2.IUserRepository
	programRepository            repositories2.IProgramRepository
	userProgramRepository        repositories2.IUserProgramRepository
	userExerciseRecordRepository repositories2.IUserExerciseRecordRepository
	exerciseRepository           repositories2.IExerciseRepository
}

func NewBot(deps botDependencies) *bot {
	b := &bot{
		shutdownWaitGroup:            deps.ShutdownWaitGroup,
		shutdownContext:              deps.ShutdownContext,
		logger:                       deps.Logger,
		config:                       deps.Config,
		senderService:                deps.SenderService,
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
		tg_bot.WithWebhookSecretToken(b.config.WebhookSecretToken()),
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

	go func() {
		bot.logger.Log("Bot started")

		if bot.config.AppEnv() == constants2.DevelopmentEnv {
			bot.bot.Start(bot.shutdownContext)
		} else {
			bot.startServer()

			bot.bot.StartWebhook(bot.shutdownContext)
		}
	}()

	for {
		select {
		case <-bot.shutdownContext.Done():
			bot.shutdown()
			return
		}
	}
}

func (bot *bot) shutdown() {
	bot.logger.Log("Shutting down gracefully...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	bot.logger.Log("Bot closed successfully")

	err := bot.server.Shutdown(shutdownCtx)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("Failed to shutdown server gracefully: %s", err))
	}

	bot.logger.Log("Server closed successfully")
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

	bot.registerCommand(constants2.CommandStart, bot.commandsHandler.Start, bot.commandMiddlewares())

	bot.registerCallbackQueryByPrefix(callback_data.MainPrefix, bot.mainHandler.Handle, bot.mainMiddlewares())

	bot.registerCallbackQueryByPrefix(callback_data.RegisterPrefix, bot.registerHandler.Handle, bot.registerMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.UserPrefix, bot.userHandler.Handle, bot.userMiddlewares())

	bot.registerCallbackQueryByPrefix(callback_data.ProgramPrefix, bot.programHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.ExercisePrefix, bot.exerciseHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.PendingUsersPrefix, bot.pendingUsersHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.BackPrefix, bot.backHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.ClientPrefix, bot.clientHandler.Handle, bot.adminMiddlewares())
}

func (bot *bot) startServer() {
	tlsConfig := &tls.Config{
		ClientAuth: tls.NoClientCert,
		MinVersion: tls.VersionTLS11,
	}

	port := bot.config.HttpPort()

	bot.server = http.Server{
		Addr:         port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      bot.bot.WebhookHandler(),
	}

	go func() {
		bot.logger.Log(fmt.Sprintf("Starting https server on port %s", port))
		if err := bot.server.ListenAndServeTLS(bot.config.SSLCertPath(), bot.config.SSLKeyPath()); err != nil && err != http.ErrServerClosed {
			bot.logger.Fatal(fmt.Sprintf("Failed to start https server: %s", err))
		}
	}()
}
