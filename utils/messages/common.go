package messages

func PressStartMessage() string {
	return "Введи /start, щоб почати роботу\\."
}

func ErrorMessage() string {
	return "Виникла помилка\\. Спробуйте пізніше ще раз\\."
}

func ParamsErrorMessage(err error) string {
	return "Виникла помилка\\. Спробуйте пізніше ще раз\\. Повідомлення для розробників: " + err.Error()
}

func RequestTimeoutMessage() string {
	return "Час відповіді на запит вичерпано. Спробуйте ще раз."
}

func DefaultMessage() string {
	return "Я не розумію тебе\\.\n Введи /start, щоб почати роботу\\. "
}

func AdminOnlyMessage() string {
	return "Ця дія доступна тільки адміністраторам\\."
}
