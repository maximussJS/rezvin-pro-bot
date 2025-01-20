package db

import (
	"context"
	"fmt"
	"go.uber.org/dig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"rezvin-pro-bot/config"
	"rezvin-pro-bot/constants"
	internal_logger "rezvin-pro-bot/internal/logger"
	"sync"
)

type dbDependencies struct {
	dig.In

	ShutdownWaitGroup *sync.WaitGroup         `name:"ShutdownWaitGroup"`
	ShutdownContext   context.Context         `name:"ShutdownContext"`
	Logger            internal_logger.ILogger `name:"Logger"`
	Config            config.IConfig          `name:"Config"`
}

func NewDB(deps dbDependencies) *gorm.DB {
	dsn := deps.Config.PostgresDSN()

	logMode := logger.Error

	if deps.Config.AppEnv() == constants.ProductionEnv {
		logMode = logger.Error
	}

	dbInstance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		deps.Logger.Error(fmt.Sprintf("Failed to connect to database: %s", err))
		panic(err)
	}

	deps.Logger.Log("Connected to database")

	deps.ShutdownWaitGroup.Add(1)

	go func() {
		defer deps.ShutdownWaitGroup.Done()

		<-deps.ShutdownContext.Done()

		db, err := dbInstance.DB()

		if err != nil {
			deps.Logger.Error(fmt.Sprintf("Failed to get database instance: %s", err))
			return
		}

		closeErr := db.Close()

		if closeErr != nil {
			deps.Logger.Error(fmt.Sprintf("Failed to close database connection: %s", closeErr))
			return
		}

		deps.Logger.Log("Database connection closed")
	}()

	return dbInstance
}
