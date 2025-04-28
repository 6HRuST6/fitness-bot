package handlers

import (
	"fitness-bot/models"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandlePhotoViewCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	if !strings.HasPrefix(data, "view_photos_") {
		return
	}

	parts := strings.Split(data, "_")
	if len(parts) != 4 {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Неверный формат запроса."))
		return
	}

	count, err := strconv.Atoi(parts[2])
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Неверное количество фото."))
		return
	}

	userID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Неверный ID пользователя."))
		return
	}

	photos, err := models.GetUserPhotos(userID, count)
	if err != nil || len(photos) == 0 {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "📭 Фото не найдены."))
		return
	}

	for i := len(photos) - 1; i >= 0; i-- {
		p := photos[i]
		photo := tgbotapi.NewPhoto(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileID(p.FileID))
		if p.Comment != "" {
			photo.Caption = p.Comment
		}
		bot.Send(photo)
	}
}
