package messages

import (
	"fmt"
	"rezvin-pro-bot/src/globals"
)

func UserNotFoundMessage(userId int64) string {
	return fmt.Sprintf("Користувача з id %d не знайдено\\.", userId)
}

func UserNotApprovedMessage() string {
	return fmt.Sprintf("%s ще не підтвердив твою реєстрацію в базі клієнтів\\. Потрібно зачекати підтвердження\\.", globals.AdminName)
}
