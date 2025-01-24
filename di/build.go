package di

import (
	"go.uber.org/dig"
	"rezvin-pro-bot/config"
	"rezvin-pro-bot/handlers"
	"rezvin-pro-bot/handlers/callback_queries"
	"rezvin-pro-bot/internal/bot"
	"rezvin-pro-bot/internal/db"
	"rezvin-pro-bot/internal/http"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/repositories"
	"rezvin-pro-bot/router"
	"rezvin-pro-bot/services"
	"rezvin-pro-bot/utils"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	c = AppendDependenciesToContainer(c, getRequiredDependencies())
	c = AppendDependenciesToContainer(c, getRepositoriesDependencies())
	c = AppendDependenciesToContainer(c, getServicesDependencies())
	c = AppendDependenciesToContainer(c, getHandlersDependencies())
	c = AppendDependenciesToContainer(c, getControllersDependencies())
	c = AppendDependenciesToContainer(c, getBotDependencies())
	c = AppendDependenciesToContainer(c, getHttpServerDependencies())

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
			Constructor: repositories.NewUserRepository,
			Interface:   new(repositories.IUserRepository),
			Token:       "UserRepository",
		},
		{
			Constructor: repositories.NewProgramRepository,
			Interface:   new(repositories.IProgramRepository),
			Token:       "ProgramRepository",
		},
		{
			Constructor: repositories.NewExerciseRepository,
			Interface:   new(repositories.IExerciseRepository),
			Token:       "ExerciseRepository",
		},
		{
			Constructor: repositories.NewUserProgramRepository,
			Interface:   new(repositories.IUserProgramRepository),
			Token:       "UserProgramRepository",
		},
		{
			Constructor: repositories.NewUserExerciseRecordRepository,
			Interface:   new(repositories.IUserExerciseRecordRepository),
			Token:       "UserExerciseRecordRepository",
		},
	}
}

func getServicesDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: services.NewConversationService,
			Interface:   new(services.IConversationService),
			Token:       "ConversationService",
		},
	}
}

func getHandlersDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: handlers.NewDefaultHandler,
			Interface:   new(handlers.IDefaultHandler),
			Token:       "DefaultHandler",
		},
		{
			Constructor: handlers.NewCommandHandler,
			Interface:   new(handlers.ICommandHandler),
			Token:       "CommandHandler",
		},
		{
			Constructor: callback_queries.NewRegisterHandler,
			Interface:   new(callback_queries.IRegisterHandler),
			Token:       "RegisterHandler",
		},
		{
			Constructor: callback_queries.NewProgramHandler,
			Interface:   new(callback_queries.IProgramHandler),
			Token:       "ProgramHandler",
		},
		{
			Constructor: callback_queries.NewExerciseHandler,
			Interface:   new(callback_queries.IExerciseHandler),
			Token:       "ExerciseHandler",
		},
		{
			Constructor: callback_queries.NewPendingUsersHandler,
			Interface:   new(callback_queries.IPendingUsersHandler),
			Token:       "PendingUsersHandler",
		},
		{
			Constructor: callback_queries.NewBackHandler,
			Interface:   new(callback_queries.IBackHandler),
			Token:       "BackHandler",
		},
		{
			Constructor: callback_queries.NewClientHandler,
			Interface:   new(callback_queries.IClientHandler),
			Token:       "ClientHandler",
		},
		{
			Constructor: callback_queries.NewUserHandler,
			Interface:   new(callback_queries.IUserHandler),
			Token:       "UserHandler",
		},
		{
			Constructor: callback_queries.NewMainHandler,
			Interface:   new(callback_queries.IMainHandler),
			Token:       "MainHandler",
		},
	}
}

func getControllersDependencies() []Dependency {
	return []Dependency{}
}

func getHttpServerDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: router.NewRouter,
			Interface:   new(router.IRouter),
			Token:       "Router",
		},
		{
			Constructor: http.NewHttpServer,
			Interface:   new(http.IHttpServer),
			Token:       "HttpServer",
		},
	}
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
