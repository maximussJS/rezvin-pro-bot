package messages

import (
	"fmt"
	"rezvin-pro-bot/src/globals"
	"rezvin-pro-bot/src/utils"
)

func NoUserProgramsMessage() string {
	return fmt.Sprintf("Програм не знайдено для тебе\\. %s ще не призначив тобі жодної програми\\.", globals.AdminName)
}

func SelectUserProgramMessage() string {
	return "Вибери програму\\."
}

func UserProgramNotAssignedMessage(programName string) string {
	return fmt.Sprintf("Програму %s не тобі не призначена\\.", utils.EscapeMarkdown(programName))
}

func SelectUserProgramOptionMessage(programName string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для програми \"*%s*\"\\:", utils.EscapeMarkdown(programName))
}

func NoRecordsForUserProgramMessage(programName string) string {
	return fmt.Sprintf("Записів не знайдено для програми \"*%s*\"\\", utils.EscapeMarkdown(programName))
}

func UserProgramAssignedMessage(programName string) string {
	return fmt.Sprintf("%s призначив тобі нову програму \"*%s*\"\\.", globals.AdminName, utils.EscapeMarkdown(programName))
}

func UserProgramUnassignedMessage(programName string) string {
	return fmt.Sprintf("%s відмінив тобі програму \"*%s*\"\\.", globals.AdminName, utils.EscapeMarkdown(programName))
}
