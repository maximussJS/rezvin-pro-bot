package bot

import (
	tg_bot "github.com/go-telegram/bot"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/constants/callback_data"
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

	bot.registerCallbackQueryByPrefix(callback_data.MainPrefix, bot.mainHandler.Handle, bot.mainMiddlewares())

	bot.registerCallbackQueryByPrefix(callback_data.RegisterPrefix, bot.registerHandler.Handle, bot.emptyMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.UserPrefix, bot.userHandler.Handle, bot.userMiddlewares())

	bot.registerCallbackQueryByPrefix(callback_data.ProgramPrefix, bot.programHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.ExercisePrefix, bot.exerciseHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.PendingUsersPrefix, bot.pendingUsersHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.BackPrefix, bot.backHandler.Handle, bot.adminMiddlewares())
	bot.registerCallbackQueryByPrefix(callback_data.ClientPrefix, bot.clientHandler.Handle, bot.adminMiddlewares())
}
