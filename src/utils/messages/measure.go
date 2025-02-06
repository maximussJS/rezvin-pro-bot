package messages

import (
	"fmt"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/utils"
)

func MeasureMenuMessage() string {
	return "Вибери одну з наступних дій для замірів\\:\n"
}

func EnterMeasureNameMessage() string {
	return "Введи назву заміру\\."
}

func EnterMeasureUnitsMessage(measureName string) string {
	return fmt.Sprintf("Введи одиниці виміру для заміру \"*%s*\"\\. Наприклад \\: мм, см, м, кг, г, л, мл, шт\\.", utils.EscapeMarkdown(measureName))
}

func MeasureNameAlreadyExistsMessage(measureName string) string {
	return fmt.Sprintf("Замір з назвою \"*%s*\" вже існує\\. Cпробуй заново", utils.EscapeMarkdown(measureName))
}

func MeasureSuccessfullyAddedMessage(measureName, units string) string {
	return fmt.Sprintf("Замір \"*%s*\" з одиницями виміру \"*%s*\" успішно доданий\\.", utils.EscapeMarkdown(measureName), utils.EscapeMarkdown(units))
}

func MeasureNotFoundMessage(id uint) string {
	return fmt.Sprintf("Замір з id %d не знайдено\\.", id)
}

func MeasuresNotFoundMessage() string {
	return fmt.Sprintf("%s ще не додав жодного заміру\\.", constants.AdminName)
}

func SelectMeasureMessage() string {
	return "Вибери замір\\."
}

func SelectMeasureOptionMessage(measureName, units string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для заміру \"*%s*\" з одиницями виміру \"*%s*\" \\:", utils.EscapeMarkdown(measureName), utils.EscapeMarkdown(units))
}

func MeasureDeletedMessage(measureName string) string {
	return fmt.Sprintf("Замір \"*%s*\" успішно видалений\\.", utils.EscapeMarkdown(measureName))
}

func MeasureRenamed(oldMeasureName, newMeasureName string) string {
	return fmt.Sprintf(
		"Замір \"*%s*\" успішно перейменований на \"*%s*\"\\.",
		utils.EscapeMarkdown(oldMeasureName),
		utils.EscapeMarkdown(newMeasureName),
	)
}

func MeasureUnitsChanged(measureName, oldUnits, newUnits string) string {
	return fmt.Sprintf(
		"Одиниці виміру для заміру \"*%s*\" змінені з \"*%s*\" на \"*%s*\"\\.",
		utils.EscapeMarkdown(measureName),
		utils.EscapeMarkdown(oldUnits),
		utils.EscapeMarkdown(newUnits),
	)
}
