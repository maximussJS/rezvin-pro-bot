package dependency

import "rezvin-pro-bot/src/internal/bot"

func GetBotDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: bot.NewBot,
			Interface:   new(bot.IBot),
			Token:       "Bot",
		},
	}
}
