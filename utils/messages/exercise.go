package messages

import (
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"rezvin-pro-bot/models"
	"strings"
)

func EnterExerciseNameMessage() string {
	return "Введи назву вправи\\."
}

func ExerciseNameAlreadyExistsMessage(exerciseName string) string {
	return fmt.Sprintf("Вправа з назвою \"*%s*\" вже існує\\. Cпробуй заново", exerciseName)
}

func ExerciseSuccessfullyAddedMessage(exerciseName, programName string) string {
	return fmt.Sprintf("Вправа \"*%s*\" успішно додана до програми \"*%s*\" \\.", exerciseName, programName)
}

func NoExercisesMessage(programName string) string {
	return fmt.Sprintf("Вправ не знайдено в програмі \"*%s*\"\\. Додай нову вправу і повтори спробу\\.", programName)
}

func ExerciseNotFoundMessage(exerciseId uint) string {
	return fmt.Sprintf("Вправа з id \"*%s*\" не знайдена\\.", exerciseId)
}

func ExercisesMessage(programName string, exercises []models.Exercise) string {
	var sb strings.Builder

	sb.WriteString("Вправи програми \"*")
	sb.WriteString(tg_bot.EscapeMarkdown(programName))
	sb.WriteString("*\"\\:\n")

	for i, exercise := range exercises {
		sb.WriteString(fmt.Sprintf("%d\\. %s\n", i+1, exercise.Name))
	}

	return sb.String()
}

func ExerciseDeleteMessage(programName string) string {
	return fmt.Sprintf("Вибери вправу для видалення з програми \"*%s*\"\\.", programName)
}

func ExerciseSuccessfullyDeletedMessage(exerciseName, programName string) string {
	return fmt.Sprintf("Вправа \"*%s*\" успішно видалена з програми \"*%s*\"\\.", exerciseName, programName)
}
