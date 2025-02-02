package messages

import (
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/utils"
	"strings"
)

func SelectUserMeasureOptionMessage(measureName string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для заміру \"*%s*\"\\:", utils.EscapeMarkdown(measureName))
}

func NoUserMeasureResultsMessage(measureName string) string {
	return fmt.Sprintf("Результатів заміру \"*%s*\" не знайдено\\. Потрібно додати новий результат заміру спочатку\\.", utils.EscapeMarkdown(measureName))
}

func UserLastMeasureDeletedMessage(measureName string) string {
	return fmt.Sprintf("Останній замір \"*%s*\" успішно видалено\\.", utils.EscapeMarkdown(measureName))
}

func SelectUserMeasureMessage() string {
	return "Вибери замір\\."
}

func EnterUserMeasureValueMessage(measureName, units string) string {
	return fmt.Sprintf("Введи значення заміру \"*%s*\" з одиницями виміру \"*%s*\"\\.", utils.EscapeMarkdown(measureName), utils.EscapeMarkdown(units))
}

func UserMeasureAddedMessage(measureName, units string, value float64) string {
	floatStr := fmt.Sprintf("%.2f", value)
	return fmt.Sprintf("Замір \"*%s*\" зі значенням \"*%s*\" \"*%s*\" успішно додано\\.", utils.EscapeMarkdown(measureName), utils.EscapeMarkdown(floatStr), utils.EscapeMarkdown(units))
}

func UserMeasureResultMessage(measure models.Measure, userMeasures []models.UserMeasure) string {
	var sb strings.Builder

	sb.WriteString("Замір \"*")
	sb.WriteString(tg_bot.EscapeMarkdown(measure.Name))
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
