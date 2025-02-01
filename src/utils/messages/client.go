package messages

import (
	"fmt"
	"rezvin-pro-bot/src/utils"
)

func NoClientsMessage() string {
	return "Клієнтів не знайдено\\."
}

func SelectClientMessage() string {
	return "Вибери клієнта\\."
}

func SelectClientOptionMessage(name string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для клієнта \"*%s*\" \\:", utils.EscapeMarkdown(name))
}

func ClientResultNotFoundMessage(id uint) string {
	return fmt.Sprintf("Запис з id %d не знайдено\\.", id)
}

func NoRecordsForClientProgramMessage(name, programName string) string {
	return fmt.Sprintf("Записів не знайдено для програми \"*%s*\" клієнта \"*%s*\"\\", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name))
}
