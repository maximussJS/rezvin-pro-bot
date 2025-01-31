package bot

import (
	"fmt"
	"github.com/go-telegram/bot/models"
	"time"
)

func GetUserID(update *models.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}

	panic(fmt.Sprintf("unable to get user id from update: %v", update))
}

func GetMessageID(update *models.Update) int {
	if update.Message != nil {
		return update.Message.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Message.ID
	}

	panic(fmt.Sprintf("unable to get message id from update: %v", update))
}

func GetChatID(update *models.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Message.Chat.ID
	}

	panic(fmt.Sprintf("unable to get chat id from update: %v", update))
}

func GetFirstName(update *models.Update) string {
	if update.Message != nil {
		return update.Message.From.FirstName
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.FirstName
	}

	panic(fmt.Sprintf("unable to get first name from update: %v", update))
}

func GetLastName(update *models.Update) string {
	if update.Message != nil {
		return update.Message.From.LastName
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.LastName
	}

	panic(fmt.Sprintf("unable to get last name from update: %v", update))
}

func GetUsername(update *models.Update) string {
	if update.Message != nil {
		return update.Message.From.Username
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.Username
	}

	panic(fmt.Sprintf("unable to get username from update: %v", update))
}

func GetUpdateTimestamp(update *models.Update) time.Time {
	if update.Message != nil {
		return time.Unix(int64(update.Message.Date), 0)
	}

	if update.CallbackQuery != nil {
		return time.Unix(int64(update.CallbackQuery.Message.Message.Date), 0)
	}

	panic(fmt.Sprintf("unable to get update timestamp from update: %v", update))
}
