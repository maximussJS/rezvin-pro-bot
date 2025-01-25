package bot

import (
	"fmt"
	"github.com/go-telegram/bot/models"
	"rezvin-pro-bot/src/utils"
	"strconv"
	"strings"
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

func GetProgramId(update *models.Update) uint {
	if update.CallbackQuery == nil {
		panic(fmt.Sprintf("unable to get program id from update: %v", update))
	}

	value := strings.Split(update.CallbackQuery.Data, ":")[1]

	valueInt, err := strconv.Atoi(value)

	utils.PanicIfError(err)

	return uint(valueInt)
}

func GetSelectedUserId(update *models.Update) int64 {
	if update.CallbackQuery == nil {
		panic(fmt.Sprintf("unable to get program id from update: %v", update))
	}

	value := strings.Split(update.CallbackQuery.Data, ":")[1]

	valueInt, err := strconv.Atoi(value)

	utils.PanicIfError(err)

	return int64(valueInt)
}

func GetClientProgramId(update *models.Update) uint {
	if update.CallbackQuery == nil {
		panic(fmt.Sprintf("unable to get program id from update: %v", update))
	}

	value := strings.Split(update.CallbackQuery.Data, ":")[2]

	valueInt, err := strconv.Atoi(value)

	utils.PanicIfError(err)

	return uint(valueInt)
}

func GetExerciseId(update *models.Update) uint {
	if update.CallbackQuery == nil {
		panic(fmt.Sprintf("unable to get exercise id from update: %v", update))
	}

	value := strings.Split(update.CallbackQuery.Data, ":")[2]

	valueInt, err := strconv.Atoi(value)

	utils.PanicIfError(err)

	return uint(valueInt)
}
