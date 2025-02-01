package dependency

import "rezvin-pro-bot/src/services"

func GetServicesDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: services.NewConversationService,
			Interface:   new(services.IConversationService),
			Token:       "ConversationService",
		},
		{
			Constructor: services.NewSenderService,
			Interface:   new(services.ISenderService),
			Token:       "SenderService",
		},
		{
			Constructor: services.NewLockService,
			Interface:   new(services.ILockService),
			Token:       "LockService",
		},
		{
			Constructor: services.NewShutdownService,
			Interface:   new(services.IShutdownService),
			Token:       "ShutdownService",
		},
	}
}
