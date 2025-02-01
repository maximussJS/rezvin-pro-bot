package bot

import (
	tg_bot "github.com/go-telegram/bot"
	"rezvin-pro-bot/src/constants"
)

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

	bot.registerCommand(constants.CommandStart, bot.commandsHandler.Start, bot.emptyMiddlewares())

	bot.registerCallbackQueryByPrefix(constants.MainPrefix, bot.mainHandler.Handle, bot.mainMiddlewares())

	bot.registerCallbackQueryByPrefix(constants.RegisterPrefix, bot.registerHandler.Handle, bot.emptyMiddlewares())
	bot.registerCallbackQueryByPrefix(constants.UserProgramPrefix, bot.userProgramHandler.Handle, bot.userMiddlewares())
	bot.registerCallbackQueryByPrefix(constants.UserResultPrefix, bot.userResultHandler.Handle, bot.userMiddlewares())

	bot.registerCallbackQueryByPrefix(constants.ProgramPrefix, bot.programHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(constants.ExercisePrefix, bot.exerciseHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(constants.PendingUsersPrefix, bot.pendingUsersHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(constants.BackPrefix, bot.backHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(constants.ClientPrefix, bot.clientHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(constants.ClientProgramPrefix, bot.clientProgramHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(constants.ClientResultPrefix, bot.clientResultHandler.Handle, bot.adminMiddlewares())
}
