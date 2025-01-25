package messages

import (
	"fmt"
	"rezvin-pro-bot/src/models"
	"strings"
)

func NoClientsMessage() string {
	return "Клієнтів не знайдено\\."
}

func SelectClientMessage() string {
	return "Вибери клієнта\\."
}

func SelectClientOptionMessage(name string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для клієнта \"*%s*\" \\:", name)
}

func NoClientProgramsMessage(name string) string {
	return fmt.Sprintf("Програм не знайдено для клієнта \"*%s*\"\\. Додай програму клієнту спочатку\\.", name)
}

func ClientProgramNotFoundMessage(id uint) string {
	return fmt.Sprintf("Програму з id %d не знайдено\\.", id)
}

func NoProgramsForClientMessage(name string) string {
	return fmt.Sprintf("Програм не знайдено для клієнта \"*%s*\"\\. Схоже користувач має всі програми\\.", name)
}

func SelectClientProgramMessage(name string) string {
	return fmt.Sprintf("Вибери програму клієнта \"*%s*\"\\.", name)
}

func SelectClientProgramOptionMessage(name, programName string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для програми \"*%s*\" клієнта \"*%s*\" \\:", programName, name)
}

func ClientProgramNotAssignedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" не призначена клієнту \"*%s*\"\\.", programName, name)
}

func ClientProgramAlreadyAssignedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" вже призначена клієнту \"*%s*\"\\. Спробуй іншу програму\\.", programName, name)
}

func ClientProgramAssignedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно призначена клієнту \"*%s*\"\\.", programName, name)
}

func ClientExerciseRecordNotFoundMessage(id uint) string {
	return fmt.Sprintf("Запис з id %d не знайдено\\.", id)
}

func ClientProgramDeletedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно видалена у клієнта \"*%s*\"\\.", programName, name)
}

func NoRecordsForClientProgramMessage(name, programName string) string {
	return fmt.Sprintf("Записів не знайдено для програми \"*%s*\" клієнта \"*%s*\"\\", programName, name)
}

func ClientProgramResultsMessage(name, programName string, records []models.UserExerciseRecord) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Результати програми \"*%s*\" клієнта \"*%s*\"\\:", programName, name))

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

func ClientProgramResultsModifyMessage(name, programName string) string {
	return fmt.Sprintf("Вибери запис для редагування результатів програми \"*%s*\" клієнта \"*%s*\"\\:", programName, name)
}

func EnterClientResultMessage(name, exerciseName string) string {
	return fmt.Sprintf("Введи результат для вправи \"*%s*\" клієнта \"*%s*\"\\:", exerciseName, name)
}

func ClientProgramResultModifiedMessage(name, exerciseName string) string {
	return fmt.Sprintf("Результати вправи \"*%s*\" клієнта \"*%s*\" успішно змінено\\.", exerciseName, name)
}
