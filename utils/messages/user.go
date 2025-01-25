package messages

import (
	"fmt"
	"rezvin-pro-bot/models"
	"strings"
)

func UserNotFoundMessage(userId int64) string {
	return "Користувача з id " + string(userId) + " не знайдено\\."
}

func UserNotApprovedMessage() string {
	return "Роман ще не підтвердив твою реєстрацію в базі клієнтів\\. Потрібно зачекати підтвердження\\."
}

func NoUserProgramsMessage() string {
	return "Програм не знайдено для тебе\\. Роман ще не призначив тобі жодної програми\\."
}

func SelectUserProgramMessage() string {
	return "Вибери програму\\."
}

func UserProgramNotAssignedMessage(programName string) string {
	return fmt.Sprintf("Програму %s не тобі не призначена\\.", programName)
}

func SelectUserProgramOptionMessage(programName string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для програми \"*%s*\"\\:", programName)
}

func NoRecordsForUserProgramMessage(programName string) string {
	return fmt.Sprintf("Записів не знайдено для програми \"*%s*\"\\", programName)
}

func UserProgramResultsMessage(programName string, records []models.UserExerciseRecord) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Результати програми \"*%s*\"\\:", programName))

	groupByName := make(map[string][]models.UserExerciseRecord)

	for _, record := range records {
		groupByName[record.Exercise.Name] = append(groupByName[record.Exercise.Name], record)
	}

	for name, records := range groupByName {
		sb.WriteString(fmt.Sprintf("\n\n*%s*\\:", name))

		for _, record := range records {
			sb.WriteString(fmt.Sprintf("\n %d повторень \\- %d кг", record.Reps, record.Weight))
		}
	}

	return sb.String()
}

func UserProgramResultsModifyMessage(programName string) string {
	return fmt.Sprintf("Вибери запис для редагування результатів програми \"*%s*\"\\:", programName)
}

func EnterUserResultMessage(exerciseName string) string {
	return fmt.Sprintf("Введи результат для вправи \"*%s*\"\\:", exerciseName)
}

func UserProgramResultModifiedMessage(exerciseName string) string {
	return fmt.Sprintf("Результати вправи \"*%s*\" успішно змінено\\.", exerciseName)
}

func UserProgramAssignedMessage(programName string) string {
	return fmt.Sprintf("Роман призначив тобі нову програму \"*%s*\"\\.", programName)
}

func UserProgramUnassignedMessage(programName string) string {
	return fmt.Sprintf("Роман відмінив тобі програму \"*%s*\"\\.", programName)
}
