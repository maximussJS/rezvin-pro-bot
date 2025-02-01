package db

import (
	"context"
	"fmt"
	"go.uber.org/dig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"rezvin-pro-bot/src/config"
	"rezvin-pro-bot/src/constants"
	internal_logger "rezvin-pro-bot/src/internal/logger"
)

type IDatabase interface {
	Shutdown(ctx context.Context) error
	GetInstance() *gorm.DB
}

type databaseDependencies struct {
	dig.In

	Logger internal_logger.ILogger `name:"Logger"`
	Config config.IConfig          `name:"Config"`
}

type database struct {
	config   config.IConfig
	logger   internal_logger.ILogger
	instance *gorm.DB
}

func NewDatabase(deps databaseDependencies) *database {
	db := &database{
		config: deps.Config,
		logger: deps.Logger,
	}

	dsn := db.config.PostgresDSN()

	logMode := logger.Error

	if db.config.AppEnv() == constants.ProductionEnv {
		logMode = logger.Error
	}

	dbInstance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		db.logger.Error(fmt.Sprintf("Failed to connect to database: %s", err))
		panic(err)
	}

	db.logger.Log("Connected to database")

	db.instance = dbInstance

	return db
}

func (db *database) GetInstance() *gorm.DB {
	return db.instance
}

func (db *database) Shutdown(_ context.Context) error {
	dbInstance, err := db.instance.DB()

	if err != nil {
		return fmt.Errorf("failed to get database instance: %s", err)
	}

	closeErr := dbInstance.Close()

	if closeErr != nil {
		return fmt.Errorf("failed to close database connection: %s", closeErr)
	}

	return nil
}
