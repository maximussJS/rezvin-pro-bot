package services

import (
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"strings"
)

type ITextService interface {
	UserRegisterMessage(firstName, lastName string) string
	UserMenuMessage(firstName, lastName string) string
	DefaultMessage() string
	ErrorMessage() string
	UnapprovedUserExistsMessage() string
	ApprovedUserExistsMessage() string
	UserRegisterSuccessMessage() string
	RequestTimeoutMessage() string
	PressStartMessage() string
	AdminStartMessage() string
	AdminProgramMenuMessage() string
	AdminOnlyMessage() string

	EnterProgramNameMessage() string
	ProgramNameAlreadyExistsMessage(programName string) string
	ProgramSuccessfullyAddedMessage(programName string) string
	ProgramSuccessfullyRenamedMessage(oldProgramName, programName string) string
	SelectProgramMessage() string
	SelectProgramOptionMessage() string

	NoProgramsMessage() string

	ProgramSuccessfullyDeletedMessage(programName string) string
}

type textService struct{}

func NewTextService() *textService {
	return &textService{}
}

func (s *textService) UserRegisterMessage(firstName, lastName string) string {
	var sb strings.Builder

	sb.WriteString("Привіт, *")
	sb.WriteString(tg_bot.EscapeMarkdown(fmt.Sprintf("%s %s", firstName, lastName)))
	sb.WriteString("*")
	sb.WriteString("\\!\n")
	sb.WriteString("Тебе не має в базі клієнтів\\. \n")
	sb.WriteString("Натисни \"📲 Реєстрація\", щоб зареєструватися\\.\n")
	sb.WriteString("Необхідно буде зачекати на моє підтвердження\\.\n")
	sb.WriteString("Після цього ти зможеш користуватися цим ботом\\.\n")

	return sb.String()
}

func (s *textService) UserMenuMessage(firstName, lastName string) string {
	var sb strings.Builder

	sb.WriteString("Привіт, *")
	sb.WriteString(tg_bot.EscapeMarkdown(fmt.Sprintf("%s %s", firstName, lastName)))
	sb.WriteString("*")
	sb.WriteString("\\!\n")
	sb.WriteString("Натисни \"🚀 Переглянути результати\", щоб переглянути свої результати\\.\n")
	sb.WriteString("Натисни \"✍️ Внести результати\", щоб внести свої результати\\.\n")

	return sb.String()
}

func (s *textService) PressStartMessage() string {
	return "Введи на /start, щоб почати роботу\\."
}

func (s *textService) DefaultMessage() string {
	return "Я не розумію тебе\\.\n Введи на /start, щоб почати роботу\\. "
}

func (s *textService) ErrorMessage() string {
	return "Виникла помилка\\. Спробуйте пізніше ще раз\\."
}

func (s *textService) UnapprovedUserExistsMessage() string {
	return "Ти вже зареєстрований в базі клієнтів\\. Але Роман ще не підключив тебе до бота\\. Чекай на підтвердження\\."
}

func (s *textService) ApprovedUserExistsMessage() string {
	return "Ти вже зареєстрований в базі клієнтів\\. Введи на /start, щоб почати роботу\\."
}

func (s *textService) UserRegisterSuccessMessage() string {
	return "Ти успішно зареєстрований в базі клієнтів\\. Чекай на підтвердження від Романа\\."
}

func (s *textService) RequestTimeoutMessage() string {
	return "Час відповіді на запит вичерпано. Спробуйте ще раз."
}

func (s *textService) AdminStartMessage() string {
	var sb strings.Builder

	sb.WriteString("Привіт, *Роман*\\!\n")
	sb.WriteString("Вибери одну з наступних дій\\:\n")

	return sb.String()
}

func (s *textService) AdminProgramMenuMessage() string {
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

func (s *textService) SelectProgramOptionMessage() string {
	return "Вибери одну з наступних дій для програм\\:"
}

func (s *textService) ProgramSuccessfullyRenamedMessage(oldProgramName, programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно перейменована на \"*%s*\" \\.", oldProgramName, programName)
}

func (s *textService) ProgramSuccessfullyDeletedMessage(programName string) string {
	return fmt.Sprintf("Програма \"*%s*\" успішно видалена\\.", programName)
}
