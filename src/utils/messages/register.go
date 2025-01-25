package messages

import (
	"fmt"
	tg_bot "github.com/go-telegram/bot"
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
	return "Ти вже зареєстрований в базі клієнтів\\. Але Роман ще не підключив тебе до бота\\. Чекай на підтвердження\\."
}

func AlreadyApprovedRegister() string {
	return "Ти вже зареєстрований в базі клієнтів\\. Введи /start, щоб почати роботу\\."
}

func SuccessRegister() string {
	return "Ти успішно зареєстрований в базі клієнтів\\. Чекай на підтвердження від Романа\\."
}

func NewRegister(name string) string {
	return fmt.Sprintf("Новий користувач \"*%s*\" чекає на підтверження\\.", name)
}
