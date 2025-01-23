package messages

func UserNotFoundMessage(userId int64) string {
	return "Користувача з id " + string(userId) + " не знайдено\\."
}
