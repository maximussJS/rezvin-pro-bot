package messages

import (
	"fmt"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/utils"
	"sort"
	"strings"
)

func ClientProgramResultsMessage(name, programName string, records []models.UserResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Результати програми \"*%s*\" клієнта \"*%s*\"\\:", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name)))

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

func ClientProgramResultsModifyMessage(name, programName string) string {
	return fmt.Sprintf("Вибери запис для редагування результатів програми \"*%s*\" клієнта \"*%s*\"\\:", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name))
}

func ClientProgramResultsSelectExerciseMessage(name, programName string) string {
	return fmt.Sprintf("Вибери вправу з програми \"*%s*\" клієнта \"*%s*\"\\, які ти хочеш відредагувати:", utils.EscapeMarkdown(programName), utils.EscapeMarkdown(name))
}

func EnterClientResultMessage(name, exerciseName string) string {
	return fmt.Sprintf("Введи результат для вправи \"*%s*\" клієнта \"*%s*\"\\:", utils.EscapeMarkdown(exerciseName), utils.EscapeMarkdown(name))
}

func ClientProgramResultModifiedMessage(name, exerciseName string, reps uint) string {
	return fmt.Sprintf("Результати вправи \"*%s*\" на %d повторень клієнта \"*%s*\" успішно змінено\\.", utils.EscapeMarkdown(exerciseName), reps, utils.EscapeMarkdown(name))
}

func ClientProgramResultExerciseSelectedMessage(name, exerciseName string) string {
	return fmt.Sprintf("Вибери кількість повторень вправи \"*%s*\" ,результат яких потрібно змінити для клієнта \"*%s*\"\\.", utils.EscapeMarkdown(exerciseName), utils.EscapeMarkdown(name))
}
