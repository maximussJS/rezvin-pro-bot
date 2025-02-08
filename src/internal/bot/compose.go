package bot

import tg_bot "github.com/go-telegram/bot"

func (bot *bot) adminMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{
		bot.skipIfConversationExistsMiddleware,
		bot.answerCallbackQueryMiddleware,
		bot.isRegisteredMiddleware,
		bot.isAdminMiddleware,
		bot.parseParamsMiddleware,
		bot.validateParamsMiddleware,
	}
}

func (bot *bot) userMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{
		bot.skipIfConversationExistsMiddleware,
		bot.answerCallbackQueryMiddleware,
		bot.isRegisteredMiddleware,
		bot.isApprovedMiddleware,
		bot.parseParamsMiddleware,
		bot.validateParamsMiddleware,
	}
}

func (bot *bot) mainMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{
		bot.skipIfConversationExistsMiddleware,
		bot.answerCallbackQueryMiddleware,
		bot.isRegisteredMiddleware,
	}
}

func (bot *bot) defaultMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{
		bot.skipOtherTypesMiddleware,
		bot.timeoutMiddleware,
		bot.panicRecoveryMiddleware,
		bot.chatIdMiddleware,
		bot.forbidParallel,
	}
}

func (bot *bot) emptyMiddlewares() []tg_bot.Middleware {
	return []tg_bot.Middleware{}
}
