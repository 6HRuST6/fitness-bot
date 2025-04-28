package core

import (
	"fitness-bot/handlers"
	"fitness-bot/models"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handlePhoto(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	fileID := update.Message.Photo[len(update.Message.Photo)-1].FileID

	// Сохраняем фото в базу
	err := models.SaveUserPhoto(userID, fileID)
	if err != nil {
		log.Println("Ошибка сохранения фото:", err)
		msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при сохранении фото.")
		Bot.Send(msg)
		return
	}

	// Сохраняем, что ждём комментарий к этому фото
	handlers.MarkPendingPhoto(chatID, fileID)

	// Просим пользователя написать комментарий
	msg := tgbotapi.NewMessage(chatID, "✍️ Напиши комментарий к фото. Например: завтрак, обед, самочувствие и т.д.")
	Bot.Send(msg)
}
