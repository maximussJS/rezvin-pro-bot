package di

import (
	"go.uber.org/dig"
	"rezvin-pro-bot/src/config"
	handlers2 "rezvin-pro-bot/src/handlers"
	callback_queries2 "rezvin-pro-bot/src/handlers/callback_queries"
	"rezvin-pro-bot/src/internal/bot"
	"rezvin-pro-bot/src/internal/db"
	"rezvin-pro-bot/src/internal/logger"
	repositories2 "rezvin-pro-bot/src/repositories"
	services2 "rezvin-pro-bot/src/services"
	"rezvin-pro-bot/src/utils"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	c = AppendDependenciesToContainer(c, getRequiredDependencies())
	c = AppendDependenciesToContainer(c, getRepositoriesDependencies())
	c = AppendDependenciesToContainer(c, getServicesDependencies())
	c = AppendDependenciesToContainer(c, getHandlersDependencies())
	c = AppendDependenciesToContainer(c, getControllersDependencies())
	c = AppendDependenciesToContainer(c, getBotDependencies())

	return c
}

func AppendDependenciesToContainer(container *dig.Container, dependencies []Dependency) *dig.Container {
	for _, dep := range dependencies {
		mustProvideDependency(container, dep)
	}

	return container
}

func mustProvideDependency(container *dig.Container, dependency Dependency) {
	if dependency.Interface == nil {
		utils.PanicIfError(container.Provide(dependency.Constructor, dig.Name(dependency.Token)))
		return
	}

	utils.PanicIfError(container.Provide(
		dependency.Constructor,
		dig.As(dependency.Interface),
		dig.Name(dependency.Token),
	))
}

func getRequiredDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: logger.NewLogger,
			Interface:   new(logger.ILogger),
			Token:       "Logger",
		},
		{
			Constructor: config.NewConfig,
			Interface:   new(config.IConfig),
			Token:       "Config",
		},
		{
			Constructor: db.NewDB,
			Interface:   nil,
			Token:       "DB",
		},
	}
}

func getRepositoriesDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: repositories2.NewUserRepository,
			Interface:   new(repositories2.IUserRepository),
			Token:       "UserRepository",
		},
		{
			Constructor: repositories2.NewProgramRepository,
			Interface:   new(repositories2.IProgramRepository),
			Token:       "ProgramRepository",
		},
		{
			Constructor: repositories2.NewExerciseRepository,
			Interface:   new(repositories2.IExerciseRepository),
			Token:       "ExerciseRepository",
		},
		{
			Constructor: repositories2.NewUserProgramRepository,
			Interface:   new(repositories2.IUserProgramRepository),
			Token:       "UserProgramRepository",
		},
		{
			Constructor: repositories2.NewUserExerciseRecordRepository,
			Interface:   new(repositories2.IUserExerciseRecordRepository),
			Token:       "UserExerciseRecordRepository",
		},
		{
			Constructor: repositories2.NewLastUserMessageRepository,
			Interface:   new(repositories2.ILastUserMessageRepository),
			Token:       "LastUserMessageRepository",
		},
	}
}

func getServicesDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: services2.NewConversationService,
			Interface:   new(services2.IConversationService),
			Token:       "ConversationService",
		},
		{
			Constructor: services2.NewSenderService,
			Interface:   new(services2.ISenderService),
			Token:       "SenderService",
		},
	}
}

func getHandlersDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: handlers2.NewDefaultHandler,
			Interface:   new(handlers2.IDefaultHandler),
			Token:       "DefaultHandler",
		},
		{
			Constructor: handlers2.NewCommandHandler,
			Interface:   new(handlers2.ICommandHandler),
			Token:       "CommandHandler",
		},
		{
			Constructor: callback_queries2.NewRegisterHandler,
			Interface:   new(callback_queries2.IRegisterHandler),
			Token:       "RegisterHandler",
		},
		{
			Constructor: callback_queries2.NewProgramHandler,
			Interface:   new(callback_queries2.IProgramHandler),
			Token:       "ProgramHandler",
		},
		{
			Constructor: callback_queries2.NewExerciseHandler,
			Interface:   new(callback_queries2.IExerciseHandler),
			Token:       "ExerciseHandler",
		},
		{
			Constructor: callback_queries2.NewPendingUsersHandler,
			Interface:   new(callback_queries2.IPendingUsersHandler),
			Token:       "PendingUsersHandler",
		},
		{
			Constructor: callback_queries2.NewBackHandler,
			Interface:   new(callback_queries2.IBackHandler),
			Token:       "BackHandler",
		},
		{
			Constructor: callback_queries2.NewClientHandler,
			Interface:   new(callback_queries2.IClientHandler),
			Token:       "ClientHandler",
		},
		{
			Constructor: callback_queries2.NewUserHandler,
			Interface:   new(callback_queries2.IUserHandler),
			Token:       "UserHandler",
		},
		{
			Constructor: callback_queries2.NewMainHandler,
			Interface:   new(callback_queries2.IMainHandler),
			Token:       "MainHandler",
		},
	}
}

func getControllersDependencies() []Dependency {
	return []Dependency{}
}

func getBotDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: bot.NewBot,
			Interface:   new(bot.IBot),
			Token:       "Bot",
		},
	}
}
