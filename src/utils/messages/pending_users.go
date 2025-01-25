package messages

import "fmt"

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
	return fmt.Sprintf("Привіт, *%s*\\! Роман підтвердив твою реєстрацію в базі клієнтів\\. Ти можеш користуватися всіма функціями бота, Введи /start, щоб почати роботу\\.\\.", name)
}

func UserDeclinedMessage(name string) string {
	return fmt.Sprintf("Привіт, *%s*\\! Роман відхилив твою реєстрацію в базі клієнтів\\. Якщо у тебе є питання, звертайся до нього\\.", name)
}

func UserApprovedForAdminMessage(name string) string {
	return fmt.Sprintf("Реєстрацію користувача \"*%s*\" підтверджено\\.", name)
}

func UserDeclinedForAdminMessage(name string) string {
	return fmt.Sprintf("Реєстрацію користувача \"*%s*\" відхилено\\.", name)
}
