package messages

import (
	"fmt"
	"rezvin-pro-bot/src/utils"
)

func NoClientMeasureMessage(name string) string {
	return fmt.Sprintf("Замірів не знайдено для клієнта \"*%s*\"\\. Додай виміри клієнту спочатку\\.", utils.EscapeMarkdown(name))
}

func SelectClientMeasureMessage(name string) string {
	return fmt.Sprintf("Вибери замір клієнта \"*%s*\"\\.", utils.EscapeMarkdown(name))
}

func ClientMeasureNotFoundMessage(id uint) string {
	return fmt.Sprintf("Замір клієнта з id %d не знайдено\\.", id)
}
