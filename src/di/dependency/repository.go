package dependency

import "rezvin-pro-bot/src/repositories"

func GetRepositoriesDependencies() []Dependency {
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
			Constructor: repositories.NewUserResultRepository,
			Interface:   new(repositories.IUserResultRepository),
			Token:       "UserResultRepository",
		},
		{
			Constructor: repositories.NewLastUserMessageRepository,
			Interface:   new(repositories.ILastUserMessageRepository),
			Token:       "LastUserMessageRepository",
		},
		{
			Constructor: repositories.NewMeasureRepository,
			Interface:   new(repositories.IMeasureRepository),
			Token:       "MeasureRepository",
		},
	}
}
