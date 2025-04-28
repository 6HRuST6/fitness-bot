package core

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ClearOldUpdates(token string) {
	// Временный бот для очистки
	tempBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Println("🧹 Очищаем старые обновления...")

	updates, _ := tempBot.GetUpdates(tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   100,
		Timeout: 0,
	})

	if len(updates) > 0 {
		lastID := updates[len(updates)-1].UpdateID
		_, _ = tempBot.GetUpdates(tgbotapi.UpdateConfig{Offset: lastID + 1})
	}

	log.Println("✅ Очистка завершена")
}
