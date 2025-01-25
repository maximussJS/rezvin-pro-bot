package config

import (
	"fmt"
	"os"
	"rezvin-pro-bot/src/utils"
	"strconv"
)

func (c *config) getRequiredString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		c.logger.Error(`Environment variable "` + key + `" not found`)
		panic(fmt.Sprintf(`Environment variable "%s" not found`, key))
	}

	return value
}

func (c *config) getOptionalString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		c.logger.Warn(`Environment variable "` + key + `" not found, used default ` + defaultValue)
		value = defaultValue
	}

	return value
}

func (c *config) getOptionalInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		c.logger.Warn(fmt.Sprintf(`Environment variable "%s" not found, used default %d`, key, defaultValue))
		return defaultValue
	}

	valueInt, err := strconv.Atoi(value)

	utils.PanicIfError(err)

	return valueInt
}

func (c *config) getOptionalInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)

	if value == "" {
		c.logger.Warn(fmt.Sprintf(`Environment variable "%s" not found, used default %d`, key, defaultValue))
		return defaultValue
	}

	valueInt, err := strconv.ParseInt(value, 10, 64)

	utils.PanicIfError(err)

	return valueInt
}

func (c *config) getOptionalBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)

	if value == "" {
		c.logger.Warn(fmt.Sprintf(`Environment variable "%s" not found, used default %t`, key, defaultValue))
		return defaultValue
	}

	valueBool, err := strconv.ParseBool(value)

	utils.PanicIfError(err)

	return valueBool
}
