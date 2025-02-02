package messages

import (
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/utils"
	"strings"
)

func NoClientMeasureResultsMessage(name, measureName string) string {
	return fmt.Sprintf("Результатів заміру \"*%s*\" клієнта \"*%s*\" не знайдено\\. Потрібно додати новий результат заміру спочатку\\.", utils.EscapeMarkdown(measureName), utils.EscapeMarkdown(name))
}

func SelectClientMeasureMessage(name string) string {
	return fmt.Sprintf("Вибери замір клієнта \"*%s*\"\\.", utils.EscapeMarkdown(name))
}

func ClientMeasureNotFoundMessage(id uint) string {
	return fmt.Sprintf("Замір клієнта з id %d не знайдено\\.", id)
}

func SelectClientMeasureOptionMessage(name, measureName string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для заміру \"*%s*\" клієнта \"*%s*\" \\:", utils.EscapeMarkdown(measureName), utils.EscapeMarkdown(name))
}

func ClientLastMeasureDeletedMessage(name, measureName string) string {
	return fmt.Sprintf("Останній замір \"*%s*\" клієнта \"*%s*\" успішно видалено\\.", utils.EscapeMarkdown(measureName), utils.EscapeMarkdown(name))
}

func EnterClientMeasureValueMessage(name, measureName, units string) string {
	return fmt.Sprintf("Введи значення заміру \"*%s*\" клієнта \"*%s*\" в %s\\.", utils.EscapeMarkdown(measureName), utils.EscapeMarkdown(name), utils.EscapeMarkdown(units))
}

func ClientMeasureAddedMessage(name, measureName, units string, value float64) string {
	floatStr := fmt.Sprintf("%.2f", value)
	return fmt.Sprintf(
		"Результат заміру \"*%s*\" клієнта \"*%s*\" зі значенням %s %s успішно додано\\.",
		utils.EscapeMarkdown(measureName),
		utils.EscapeMarkdown(name),
		utils.EscapeMarkdown(floatStr),
		utils.EscapeMarkdown(units),
	)
}

func ClientMeasureResultMessage(name string, measure models.Measure, userMeasures []models.UserMeasure) string {
	var sb strings.Builder

	sb.WriteString("Замір \"*")
	sb.WriteString(tg_bot.EscapeMarkdown(measure.Name))
	sb.WriteString("*\" клієнта \"*")
	sb.WriteString(tg_bot.EscapeMarkdown(name))
	sb.WriteString("*\"\\:\n")

	for i, userMeasure := range userMeasures {
		floatStr := fmt.Sprintf("%.2f", userMeasure.Value)

		sb.WriteString(fmt.Sprintf(
			"%d\\. %s \\- %s %s\\.\n",
			i+1,
			utils.EscapeMarkdown(userMeasure.CreatedAt.Format("2006-01-02 15:04:05")),
			utils.EscapeMarkdown(floatStr),
			utils.EscapeMarkdown(measure.Units),
		))
	}

	return sb.String()
}
