package globals

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func GetPostgresSchema() string {
	key := "POSTGRES_SCHEMA"
	godotenv.Load() // ignore error, because in deployment we pass all env variables via docker run command

	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf(`Environment variable "%s" not found`, key))
	}

	return value
}

func GetAdminName() string {
	key := "ADMIN_NAME"
	godotenv.Load() // ignore error, because in deployment we pass all env variables via docker run command

	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf(`Environment variable "%s" not found`, key))
	}

	return value
}

var AdminName = GetAdminName()
