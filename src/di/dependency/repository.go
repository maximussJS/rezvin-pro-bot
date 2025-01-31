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
			Constructor: repositories.NewUserExerciseRecordRepository,
			Interface:   new(repositories.IUserExerciseRecordRepository),
			Token:       "UserExerciseRecordRepository",
		},
		{
			Constructor: repositories.NewLastUserMessageRepository,
			Interface:   new(repositories.ILastUserMessageRepository),
			Token:       "LastUserMessageRepository",
		},
	}
}
