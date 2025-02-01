package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/internal/logger"
	"time"
)

type IConfig interface {
	AppEnv() constants.AppEnv

	BotToken() string
	WebhookSecretToken() string
	RequestTimeout() time.Duration
	AlertChatId() int64

	ErrorStackTraceSizeInKb() int

	PostgresDSN() string
	RunMigrations() bool

	HttpPort() string
	SSLCertPath() string
	SSLKeyPath() string
}

type configDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`
}

type config struct {
	logger logger.ILogger

	appEnv constants.AppEnv

	botToken                string
	webhookSecretToken      string
	requestTimeoutInSeconds int
	alertChatId             int64

	errorStackTraceSizeInKb int

	postgresDsn   string
	runMigrations bool

	httpPort    string
	sslCertPath string
	sslKeyPath  string
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
		config.webhookSecretToken = config.getRequiredString("WEBHOOK_SECRET_TOKEN")
		config.sslCertPath = config.getOptionalString("SSL_CERT_PATH", "./certs/cert.pem")
		config.sslKeyPath = config.getOptionalString("SSL_KEY_PATH", "./certs/priv.pem")
	default:
		panic(fmt.Sprintf("Invalid APP_ENV value: %s. Supported values: %s, %s", appEnv, constants.DevelopmentEnv, constants.ProductionEnv))
	}

	config.botToken = config.getRequiredString("BOT_TOKEN")
	config.postgresDsn = config.getRequiredString("POSTGRES_DSN")
	config.alertChatId = config.getRequiredInt64("ALERT_CHAT_ID")
	config.runMigrations = config.getOptionalBool("RUN_MIGRATIONS", false)
	config.requestTimeoutInSeconds = config.getOptionalInt("REQUEST_TIMEOUT_IN_SECONDS", 60)
	config.errorStackTraceSizeInKb = config.getOptionalInt("ERROR_STACK_TRACE_SIZE_IN_KB", 4)
	config.httpPort = config.getOptionalString("HTTP_PORT", ":8080")

	return config
}

func (c *config) AppEnv() constants.AppEnv {
	return c.appEnv
}

func (c *config) BotToken() string {
	return c.botToken
}

func (c *config) SSLCertPath() string {
	return c.sslCertPath
}

func (c *config) WebhookSecretToken() string {
	return c.webhookSecretToken
}

func (c *config) SSLKeyPath() string {
	return c.sslKeyPath
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

func (c *config) HttpPort() string {
	return c.httpPort
}

func (c *config) AlertChatId() int64 {
	return c.alertChatId
}
