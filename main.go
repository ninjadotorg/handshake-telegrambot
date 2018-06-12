package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if botToken == "" || chatID == 0 {
		panic(errors.New("env is invalid"))
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	joinGroupTemplate, err := makeTemplate("templates/join_group.html", "join_group.html")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.ID == chatID {
			if update.Message.NewChatMember != nil {
				newChatMember := *update.Message.NewChatMember
				messageHTML, err := makeContent(joinGroupTemplate, "join_group.html", newChatMember)
				if err != nil {
					log.Println(err)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageHTML)
				msg.ParseMode = "html"
				bot.Send(msg)
			}
		}
	}

}

func makeTemplate(filePath string, name string) (*template.Template, error) {
	t, err := template.New(name).ParseFiles(filePath)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func makeContent(t *template.Template, name string, data interface{}) (string, error) {
	var tpl bytes.Buffer
	err := t.ExecuteTemplate(&tpl, name, data)
	if err != nil {
		return "", err
	}
	content := tpl.String()
	return content, nil
}
