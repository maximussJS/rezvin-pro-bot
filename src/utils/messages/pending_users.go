package messages

import (
	"fmt"
	"rezvin-pro-bot/src/constants"
)

func NoPendingUsersMessage() string {
	return "Немає користувачів, які чекають на підтвердження\\."
}

func SelectPendingUserMessage() string {
	return "Вибери користувача для підтвердження\\."
}

func SelectPendingUserOptionMessage(name string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для користувача \"*%s*\" \\:", name)
}

func UserApprovedMessage(name string) string {
	return fmt.Sprintf("Привіт, *%s*\\! %s підтвердив твою реєстрацію в базі клієнтів\\. Ти можеш користуватися всіма функціями бота, Введи /start, щоб почати роботу\\.\\.", name, constants.AdminName)
}

func UserDeclinedMessage(name string) string {
	return fmt.Sprintf("Привіт, *%s*\\! %s відхилив твою реєстрацію в базі клієнтів\\. Якщо у тебе є питання, звертайся до нього\\.", name, constants.AdminName)
}

func UserApprovedForAdminMessage(name string) string {
	return fmt.Sprintf("Реєстрацію користувача \"*%s*\" підтверджено\\.", name)
}

func UserDeclinedForAdminMessage(name string) string {
	return fmt.Sprintf("Реєстрацію користувача \"*%s*\" відхилено\\.", name)
}
