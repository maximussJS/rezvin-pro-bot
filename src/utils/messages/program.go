package messages

import (
	"fmt"
	"rezvin-pro-bot/src/utils"
)

func ProgramMenuMessage() string {
	return "Вибери одну з наступних дій для програм\\:\n"
}

func EnterProgramNameMessage() string {
	return "Введи назву програми\\."
}

func ProgramNameAlreadyExistsMessage(programName string) string {
	return fmt.Sprintf("Програма з назвою \"*%s*\" вже існує\\. Cпробуй заново", utils.EscapeMarkdown(programName))
}

func ProgramSuccessfullyAddedMessage(programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно додана\\.", utils.EscapeMarkdown(programName))
}

func NoProgramsMessage() string {
	return "Програм не знайдено\\. Створи нову програму і повтори спробу\\."
}

func ProgramNotFoundMessage(programId uint) string {
	return fmt.Sprintf("Програму з id %d не знайдено\\.", programId)
}

func SelectProgramMessage() string {
	return "Вибери програму\\."
}

func SelectProgramOptionMessage(programName string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для програми \"*%s*\" \\:", utils.EscapeMarkdown(programName))
}

func ProgramSuccessfullyRenamedMessage(oldProgramName, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно перейменована на \"*%s*\" \\.", utils.EscapeMarkdown(oldProgramName), utils.EscapeMarkdown(programName))
}

func ProgramSuccessfullyDeletedMessage(programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно видалена\\.", utils.EscapeMarkdown(programName))
}
