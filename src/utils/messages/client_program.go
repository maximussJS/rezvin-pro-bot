package messages

import (
	"fmt"
	"rezvin-pro-bot/src/utils"
)

func NoClientProgramsMessage(name string) string {
	return fmt.Sprintf("Програм не знайдено для клієнта \"*%s*\"\\. Додай програму клієнту спочатку\\.", utils.EscapeMarkdown(name))
}

func ClientProgramNotFoundMessage(id uint) string {
	return fmt.Sprintf("Програму з id %d не знайдено\\.", id)
}

func NoProgramsForClientMessage(name string) string {
	return fmt.Sprintf("Програм не знайдено для клієнта \"*%s*\"\\. Схоже користувач має всі програми\\.", utils.EscapeMarkdown(name))
}

func SelectClientProgramMessage(name string) string {
	return fmt.Sprintf("Вибери програму клієнта \"*%s*\"\\.", utils.EscapeMarkdown(name))
}

func SelectClientProgramOptionMessage(name, programName string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для програми \"*%s*\" клієнта \"*%s*\" \\:", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name))
}

func ClientProgramNotAssignedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" не призначена клієнту \"*%s*\"\\.", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name))
}

func ClientProgramAlreadyAssignedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" вже призначена клієнту \"*%s*\"\\. Спробуй іншу програму\\.", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name))
}

func ClientProgramAssignedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно призначена клієнту \"*%s*\"\\.", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name))
}

func ClientProgramDeletedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно видалена у клієнта \"*%s*\"\\.", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name))
}
