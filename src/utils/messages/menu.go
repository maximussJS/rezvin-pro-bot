package messages

import (
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"rezvin-pro-bot/src/globals"
	"strings"
)

func UserMenuMessage(name string) string {
	var sb strings.Builder

	sb.WriteString("Привіт, *")
	sb.WriteString(tg_bot.EscapeMarkdown(fmt.Sprintf("%s", name)))
	sb.WriteString("*")
	sb.WriteString("\\!\n")
	sb.WriteString("Натисни \"🚀 Переглянути результати\", щоб переглянути свої результати\\.\n")
	sb.WriteString("Натисни \"✍️ Внести результати\", щоб внести свої результати\\.\n")

	return sb.String()
}

func AdminMainMessage() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Привіт, *%s*\\!\n", globals.AdminName))
	sb.WriteString("Вибери одну з наступних дій\\:\n")

	return sb.String()
}
