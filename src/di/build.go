package di

import (
	"go.uber.org/dig"
	"rezvin-pro-bot/src/di/dependency"
	"rezvin-pro-bot/src/utils"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	c = AppendDependenciesToContainer(c, dependency.GetRequiredDependencies())
	c = AppendDependenciesToContainer(c, dependency.GetRepositoriesDependencies())
	c = AppendDependenciesToContainer(c, dependency.GetServicesDependencies())
	c = AppendDependenciesToContainer(c, dependency.GetHandlersDependencies())
	c = AppendDependenciesToContainer(c, dependency.GetBotDependencies())

	return c
}

func AppendDependenciesToContainer(container *dig.Container, dependencies []dependency.Dependency) *dig.Container {
	for _, dep := range dependencies {
		mustProvideDependency(container, dep)
	}

	return container
}

func mustProvideDependency(container *dig.Container, dependency dependency.Dependency) {
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
