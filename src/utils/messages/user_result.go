package messages

import (
	"fmt"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/utils"
	"sort"
	"strings"
)

func UserProgramResultsMessage(programName string, records []models.UserResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Результати програми \"*%s*\"\\:", utils.EscapeMarkdown(programName)))

	groupByName := make(map[string][]models.UserResult)
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

func UserProgramResultsSelectExerciseMessage(programName string) string {
	return fmt.Sprintf("Вибери вправу з програми \"*%s*\", яку ти хочеш відредагувати:", utils.EscapeMarkdown(programName))
}

func EnterUserResultMessage(exerciseName string) string {
	return fmt.Sprintf("Введи результат для вправи \"*%s*\"\\:", utils.EscapeMarkdown(exerciseName))
}

func UserProgramResultModifiedMessage(exerciseName string, reps uint) string {
	return fmt.Sprintf("Результати вправи \"*%s*\"на %d повторень успішно змінено\\.", utils.EscapeMarkdown(exerciseName), reps)
}

func UserProgramResultExerciseSelectedMessage(exerciseName string) string {
	return fmt.Sprintf("Вибери кількість повторень вправи \"*%s*\" ,результат яких потрібно змінити \\.", utils.EscapeMarkdown(exerciseName))
}
