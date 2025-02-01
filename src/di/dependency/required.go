package dependency

import (
	"rezvin-pro-bot/src/config"
	"rezvin-pro-bot/src/internal/db"
	"rezvin-pro-bot/src/internal/logger"
)

func GetRequiredDependencies() []Dependency {
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
			Constructor: db.NewDatabase,
			Interface:   new(db.IDatabase),
			Token:       "Database",
		},
	}
}
