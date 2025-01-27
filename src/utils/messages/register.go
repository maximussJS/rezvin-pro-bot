package messages

import (
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/utils"
	"strings"
)

func NeedRegister(name string) string {
	var sb strings.Builder

	sb.WriteString("Привіт, *")
	sb.WriteString(tg_bot.EscapeMarkdown(fmt.Sprintf("%s", name)))
	sb.WriteString("*")
	sb.WriteString("\\!\n")
	sb.WriteString("Тебе не має в базі клієнтів\\. \n")
	sb.WriteString("Натисни \"📲 Реєстрація\", щоб зареєструватися\\.\n")
	sb.WriteString("Необхідно буде зачекати на моє підтвердження\\.\n")
	sb.WriteString("Після цього ти зможеш користуватися цим ботом\\.\n")

	return sb.String()
}

func AlreadyRegistered() string {
	return fmt.Sprintf("Ти вже зареєстрований в базі клієнтів\\. Але %s ще не підключив тебе до бота\\. Чекай на підтвердження\\.", constants.AdminName)
}

func AlreadyApprovedRegister() string {
	return "Ти вже зареєстрований в базі клієнтів\\. Введи /start, щоб почати роботу\\."
}

func SuccessRegister() string {
	return fmt.Sprintf("Ти успішно зареєстрований в базі клієнтів\\. Чекай на підтвердження від %s\\.", constants.AdminName)
}

func NewRegister(name string) string {
	return fmt.Sprintf("Новий користувач \"*%s*\" чекає на підтверження\\.", utils.EscapeMarkdown(name))
}
