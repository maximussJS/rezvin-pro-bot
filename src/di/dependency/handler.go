package dependency

import (
	"rezvin-pro-bot/src/handlers"
	cb_handlers "rezvin-pro-bot/src/handlers/callback_queries"
)

func GetHandlersDependencies() []Dependency {
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
			Constructor: cb_handlers.NewRegisterHandler,
			Interface:   new(cb_handlers.IRegisterHandler),
			Token:       "RegisterHandler",
		},
		{
			Constructor: cb_handlers.NewProgramHandler,
			Interface:   new(cb_handlers.IProgramHandler),
			Token:       "ProgramHandler",
		},
		{
			Constructor: cb_handlers.NewExerciseHandler,
			Interface:   new(cb_handlers.IExerciseHandler),
			Token:       "ExerciseHandler",
		},
		{
			Constructor: cb_handlers.NewPendingUsersHandler,
			Interface:   new(cb_handlers.IPendingUsersHandler),
			Token:       "PendingUsersHandler",
		},
		{
			Constructor: cb_handlers.NewBackHandler,
			Interface:   new(cb_handlers.IBackHandler),
			Token:       "BackHandler",
		},
		{
			Constructor: cb_handlers.NewClientHandler,
			Interface:   new(cb_handlers.IClientHandler),
			Token:       "ClientHandler",
		},
		{
			Constructor: cb_handlers.NewClientProgramHandler,
			Interface:   new(cb_handlers.IClientProgramHandler),
			Token:       "ClientProgramHandler",
		},
		{
			Constructor: cb_handlers.NewClientResultHandler,
			Interface:   new(cb_handlers.IClientResultHandler),
			Token:       "ClientResultHandler",
		},
		{
			Constructor: cb_handlers.NewUserResultHandler,
			Interface:   new(cb_handlers.IUserResultHandler),
			Token:       "UserResultHandler",
		},
		{
			Constructor: cb_handlers.NewUserProgramHandler,
			Interface:   new(cb_handlers.IUserProgramHandler),
			Token:       "UserProgramHandler",
		},
		{
			Constructor: cb_handlers.NewMainHandler,
			Interface:   new(cb_handlers.IMainHandler),
			Token:       "MainHandler",
		},
		{
			Constructor: cb_handlers.NewMeasureHandler,
			Interface:   new(cb_handlers.IMeasureHandler),
			Token:       "MeasureHandler",
		},
		{
			Constructor: cb_handlers.NewClientMeasureHandler,
			Interface:   new(cb_handlers.IClientMeasureHandler),
			Token:       "ClientMeasureHandler",
		},
	}
}
