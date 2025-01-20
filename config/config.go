package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
	"rezvin-pro-bot/constants"
	"rezvin-pro-bot/internal/logger"
	"time"
)

type IConfig interface {
	AppEnv() constants.AppEnv

	BotToken() string

	ErrorStackTraceSizeInKb() int

	PostgresDSN() string
	RunMigrations() bool

	HttpPort() string
	MaxJSONBodySizeInBytes() int64
	RequestTimeout() time.Duration

	AuthToken() string
}

type configDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`
}

type config struct {
	logger logger.ILogger

	appEnv constants.AppEnv

	botToken string

	errorStackTraceSizeInKb int

	postgresDsn   string
	runMigrations bool

	httpPort                string
	maxJSONBodySizeInBytes  int64
	requestTimeoutInSeconds int

	authToken string
}

func NewConfig(deps configDependencies) *config {
	_logger := deps.Logger

	godotenv.Load() // ignore error, because in deployment we pass all env variables via docker run command

	config := &config{
		logger: _logger,
	}

	appEnv := config.getRequiredString("APP_ENV")

	switch appEnv {
	case string(constants.DevelopmentEnv):
		config.appEnv = constants.DevelopmentEnv
	case string(constants.ProductionEnv):
		config.appEnv = constants.ProductionEnv
	default:
		panic(fmt.Sprintf("Invalid APP_ENV value: %s. Supported values: %s, %s", appEnv, constants.DevelopmentEnv, constants.ProductionEnv))
	}

	config.botToken = config.getRequiredString("BOT_TOKEN")
	config.authToken = config.getRequiredString("AUTH_TOKEN")
	config.postgresDsn = config.getRequiredString("POSTGRES_DSN")
	config.runMigrations = config.getOptionalBool("RUN_MIGRATIONS", false)
	config.errorStackTraceSizeInKb = config.getOptionalInt("ERROR_STACK_TRACE_SIZE_IN_KB", 4)
	config.httpPort = config.getOptionalString("HTTP_PORT", ":8080")
	config.maxJSONBodySizeInBytes = config.getOptionalInt64("MAX_JSON_BODY_SIZE_IN_BYTES", 1048576)
	config.requestTimeoutInSeconds = config.getOptionalInt("REQUEST_TIMEOUT_IN_SECONDS", 15)

	return config
}

func (c *config) AppEnv() constants.AppEnv {
	return c.appEnv
}

func (c *config) BotToken() string {
	return c.botToken
}

func (c *config) ErrorStackTraceSizeInKb() int {
	return c.errorStackTraceSizeInKb
}

func (c *config) PostgresDSN() string {
	return c.postgresDsn
}

func (c *config) RunMigrations() bool {
	return c.runMigrations
}

func (c *config) RequestTimeout() time.Duration {
	return time.Duration(c.requestTimeoutInSeconds) * time.Second
}

func (c *config) MaxJSONBodySizeInBytes() int64 {
	return c.maxJSONBodySizeInBytes
}

func (c *config) HttpPort() string {
	return c.httpPort
}

func (c *config) AuthToken() string {
	return c.authToken
}
