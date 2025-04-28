package core

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func Start(token string) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	Bot.Debug = true
	log.Printf("Бот запущен: @%s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(1)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallback(update)
			continue
		}
		if update.Message != nil {
			handleMessage(update)
		}
	}
}
