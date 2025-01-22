package services

import (
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"rezvin-pro-bot/models"
	"strings"
)

type ITextService interface {
	UserRegisterMessage(name string) string
	UserMenuMessage(name string) string
	DefaultMessage() string
	ErrorMessage() string
	UnapprovedUserExistsMessage() string
	DeclinedUserExistsMessage() string
	ApprovedUserExistsMessage() string
	UserRegisterSuccessMessage() string
	RequestTimeoutMessage() string
	PressStartMessage() string
	AdminMainMessage() string
	ProgramMenuMessage() string
	AdminOnlyMessage() string
	NewUserRegisteredMessage(name string) string

	EnterProgramNameMessage() string
	ProgramNameAlreadyExistsMessage(programName string) string
	ProgramSuccessfullyAddedMessage(programName string) string
	ProgramSuccessfullyRenamedMessage(oldProgramName, programName string) string
	SelectProgramMessage() string
	SelectProgramOptionMessage(programName string) string

	NoProgramsMessage() string

	ProgramSuccessfullyDeletedMessage(programName string) string

	EnterExerciseNameMessage() string
	ExerciseNameAlreadyExistsMessage(exerciseName string) string
	ExerciseSuccessfullyAddedMessage(exerciseName, programName string) string
	NoExercisesMessage(programName string) string
	ExercisesMessage(programName string, exercises []models.Exercise) string
	ExerciseDeleteMessage(programName string) string
	ExerciseSuccessfullyDeletedMessage(exerciseName, programName string) string

	NoPendingUsersMessage() string
	SelectPendingUserMessage() string
	SelectPendingUserOptionMessage(name string) string

	UserApprovedMessage(name string) string
	UserApprovedForAdminMessage(name string) string

	UserDeclinedMessage(name string) string
	UserDeclinedForAdminMessage(name string) string

	NoClientsMessage() string
	SelectClientMessage() string
	SelectClientOptionMessage(name string) string
	NoClientProgramsMessage(name string) string
	SelectClientProgramMessage(name string) string
	ClientProgramAssignedMessage(name, programName string) string
	NoProgramsForClientMessage(name string) string
	ClientProgramAlreadyAssignedMessage(name, programName string) string
}

type textService struct{}

func NewTextService() *textService {
	return &textService{}
}

func (s *textService) UserRegisterMessage(name string) string {
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

func (s *textService) UserMenuMessage(name string) string {
	var sb strings.Builder

	sb.WriteString("Привіт, *")
	sb.WriteString(tg_bot.EscapeMarkdown(fmt.Sprintf("%s", name)))
	sb.WriteString("*")
	sb.WriteString("\\!\n")
	sb.WriteString("Натисни \"🚀 Переглянути результати\", щоб переглянути свої результати\\.\n")
	sb.WriteString("Натисни \"✍️ Внести результати\", щоб внести свої результати\\.\n")

	return sb.String()
}

func (s *textService) PressStartMessage() string {
	return "Введи /start, щоб почати роботу\\."
}

func (s *textService) DefaultMessage() string {
	return "Я не розумію тебе\\.\n Введи /start, щоб почати роботу\\. "
}

func (s *textService) ErrorMessage() string {
	return "Виникла помилка\\. Спробуйте пізніше ще раз\\."
}

func (s *textService) UnapprovedUserExistsMessage() string {
	return "Ти вже зареєстрований в базі клієнтів\\. Але Роман ще не підключив тебе до бота\\. Чекай на підтвердження\\."
}

func (s *textService) ApprovedUserExistsMessage() string {
	return "Ти вже зареєстрований в базі клієнтів\\. Введи /start, щоб почати роботу\\."
}

func (s *textService) UserRegisterSuccessMessage() string {
	return "Ти успішно зареєстрований в базі клієнтів\\. Чекай на підтвердження від Романа\\."
}

func (s *textService) RequestTimeoutMessage() string {
	return "Час відповіді на запит вичерпано. Спробуйте ще раз."
}

func (s *textService) AdminMainMessage() string {
	var sb strings.Builder

	sb.WriteString("Привіт, *Роман*\\!\n")
	sb.WriteString("Вибери одну з наступних дій\\:\n")

	return sb.String()
}

func (s *textService) ProgramMenuMessage() string {
	return "Вибери одну з наступних дій для програм\\:\n"
}

func (s *textService) AdminOnlyMessage() string {
	return "Ця команда доступна тільки адміністраторам\\."
}

func (s *textService) EnterProgramNameMessage() string {
	return "Введи назву програми\\."
}

func (s *textService) ProgramNameAlreadyExistsMessage(programName string) string {
	return fmt.Sprintf("Програма з назвою \"*%s*\" вже існує\\. Cпробуй заново", programName)
}

func (s *textService) ProgramSuccessfullyAddedMessage(programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно додана\\.", programName)
}

func (s *textService) NoProgramsMessage() string {
	return "Програм не знайдено\\. Створи нову програму і повтори спробу\\."
}

func (s *textService) SelectProgramMessage() string {
	return "Вибери програму\\."
}

func (s *textService) SelectProgramOptionMessage(programName string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для програми \"*%s*\" \\:", programName)
}

func (s *textService) ProgramSuccessfullyRenamedMessage(oldProgramName, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно перейменована на \"*%s*\" \\.", oldProgramName, programName)
}

func (s *textService) ProgramSuccessfullyDeletedMessage(programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно видалена\\.", programName)
}

func (s *textService) EnterExerciseNameMessage() string {
	return "Введи назву вправи\\."
}

func (s *textService) ExerciseNameAlreadyExistsMessage(exerciseName string) string {
	return fmt.Sprintf("Вправа з назвою \"*%s*\" вже існує\\. Cпробуй заново", exerciseName)
}

func (s *textService) ExerciseSuccessfullyAddedMessage(exerciseName, programName string) string {
	return fmt.Sprintf("Вправа \"*%s*\" успішно додана до програми \"*%s*\" \\.", exerciseName, programName)
}

func (s *textService) NoExercisesMessage(programName string) string {
	return fmt.Sprintf("Вправ не знайдено в програмі \"*%s*\"\\. Додай нову вправу і повтори спробу\\.", programName)
}

func (s *textService) ExercisesMessage(programName string, exercises []models.Exercise) string {
	var sb strings.Builder

	sb.WriteString("Вправи програми \"*")
	sb.WriteString(tg_bot.EscapeMarkdown(programName))
	sb.WriteString("*\"\\:\n")

	for i, exercise := range exercises {
		sb.WriteString(fmt.Sprintf("%d\\. %s\n", i+1, exercise.Name))
	}

	return sb.String()
}

func (s *textService) ExerciseDeleteMessage(programName string) string {
	return fmt.Sprintf("Вибери вправу для видалення з програми \"*%s*\"\\.", programName)
}

func (s *textService) ExerciseSuccessfullyDeletedMessage(exerciseName, programName string) string {
	return fmt.Sprintf("Вправа \"*%s*\" успішно видалена з програми \"*%s*\"\\.", exerciseName, programName)
}

func (s *textService) NoPendingUsersMessage() string {
	return "Немає користувачів, які чекають на підтвердження\\."
}

func (s *textService) SelectPendingUserMessage() string {
	return "Вибери користувача для підтвердження\\."
}

func (s *textService) SelectPendingUserOptionMessage(name string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для користувача \"*%s*\" \\:", name)
}

func (s *textService) UserApprovedMessage(name string) string {
	return fmt.Sprintf("Привіт, *%s*\\! Роман підтвердив твою реєстрацію в базі клієнтів\\. Тепер ти можеш користуватися всіма функціями бота, Введи /start, щоб почати роботу\\.\\.", name)
}

func (s *textService) UserDeclinedMessage(name string) string {
	return fmt.Sprintf("Привіт, *%s*\\! Роман відхилив твою реєстрацію в базі клієнтів\\. Якщо у тебе є питання, звертайся до нього\\.", name)
}

func (s *textService) UserApprovedForAdminMessage(name string) string {
	return fmt.Sprintf("Реєстрацію користувача \"*%s*\" підтверджено\\.", name)
}

func (s *textService) DeclinedUserExistsMessage() string {
	return "Роман відхилив твою реєстрацію\\."
}

func (s *textService) UserDeclinedForAdminMessage(name string) string {
	return fmt.Sprintf("Реєстрацію користувача \"*%s*\" відхилено\\.", name)
}

func (s *textService) NewUserRegisteredMessage(name string) string {
	return fmt.Sprintf("Новий користувач \"*%s*\" чекає на підтверження\\.", name)
}

func (s *textService) NoClientsMessage() string {
	return "Клієнтів не знайдено\\."
}

func (s *textService) SelectClientMessage() string {
	return "Вибери клієнта\\."
}

func (s *textService) SelectClientOptionMessage(name string) string {
	return fmt.Sprintf("Вибери одну з наступних дій для клієнта \"*%s*\" \\:", name)
}

func (s *textService) NoClientProgramsMessage(name string) string {
	return fmt.Sprintf("Програм не знайдено для клієнта \"*%s*\"\\. Додай програму клієнту спочатку\\.", name)
}

func (s *textService) NoProgramsForClientMessage(name string) string {
	return fmt.Sprintf("Програм не знайдено для клієнта \"*%s*\"\\. Схоже користувач має всі програми\\.", name)
}

func (s *textService) SelectClientProgramMessage(name string) string {
	return fmt.Sprintf("Вибери програму клієнта \"*%s*\"\\.", name)
}

func (s *textService) ClientProgramAlreadyAssignedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" вже призначена клієнту \"*%s*\"\\. Спробуй іншу програму\\.", programName, name)
}

func (s *textService) ClientProgramAssignedMessage(name, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно призначена клієнту \"*%s*\"\\.", programName, name)
}
