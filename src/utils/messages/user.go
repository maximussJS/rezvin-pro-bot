package messages

import (
	"fmt"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/utils"
	"sort"
	"strings"
)

func UserNotFoundMessage(userId int64) string {
	return "Користувача з id " + string(userId) + " не знайдено\\."
}

func UserNotApprovedMessage() string {
	return fmt.Sprintf("%s ще не підтвердив твою реєстрацію в базі клієнтів\\. Потрібно зачекати підтвердження\\.", constants.AdminName)
}

func NoUserProgramsMessage() string {
	return fmt.Sprintf("Програм не знайдено для тебе\\. %s ще не призначив тобі жодної програми\\.", constants.AdminName)
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

func UserProgramResultsMessage(programName string, records []models.UserExerciseRecord) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Результати програми \"*%s*\"\\:", utils.EscapeMarkdown(programName)))

	groupByName := make(map[string][]models.UserExerciseRecord)
	groupById := make(map[uint]string)

	for _, record := range records {
		groupById[record.ExerciseId] = record.Exercise.Name
		groupByName[record.Exercise.Name] = append(groupByName[record.Exercise.Name], record)
	}

	sortedIds := make([]int, 0, len(groupById))

	for id := range groupById {
		sortedIds = append(sortedIds, int(id))
	}

	sort.Sort(sort.IntSlice(sortedIds))

	for _, id := range sortedIds {
		name := groupById[uint(id)]
		records := groupByName[name]

		sb.WriteString(fmt.Sprintf("\n\n*%s*\\:", utils.EscapeMarkdown(name)))

		for _, record := range records {
			sb.WriteString(fmt.Sprintf("\n %d повторень \\- %d кг", record.Reps, record.Weight))
		}
	}

	return sb.String()
}

func UserProgramResultsModifyMessage(programName string) string {
	return fmt.Sprintf("Вибери запис для редагування результатів програми \"*%s*\"\\:", utils.EscapeMarkdown(programName))
}

func EnterUserResultMessage(exerciseName string) string {
	return fmt.Sprintf("Введи результат для вправи \"*%s*\"\\:", utils.EscapeMarkdown(exerciseName))
}

func UserProgramResultModifiedMessage(exerciseName string) string {
	return fmt.Sprintf("Результати вправи \"*%s*\" успішно змінено\\.", utils.EscapeMarkdown(exerciseName))
}

func UserProgramAssignedMessage(programName string) string {
	return fmt.Sprintf("%s призначив тобі нову програму \"*%s*\"\\.", constants.AdminName, utils.EscapeMarkdown(programName))
}

func UserProgramUnassignedMessage(programName string) string {
	return fmt.Sprintf("%s відмінив тобі програму \"*%s*\"\\.", constants.AdminName, utils.EscapeMarkdown(programName))
}
