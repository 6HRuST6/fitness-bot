package handlers

import (
	"fitness-bot/models"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var pendingPhotoComments = make(map[int64]string)

func MarkPendingPhoto(userID int64, fileID string) {
	pendingPhotoComments[userID] = fileID
}

func PendingComment(userID int64) (string, bool) {
	fileID, ok := pendingPhotoComments[userID]
	return fileID, ok
}

func SavePhotoComment(bot *tgbotapi.BotAPI, update tgbotapi.Update, fileID string) {
	userID := update.Message.Chat.ID
	username := update.Message.From.UserName
	comment := update.Message.Text

	delete(pendingPhotoComments, userID)

	models.AddPhotoComment(userID, comment)

	text := fmt.Sprintf("🗒 Комментарий от @%s:\n%s", username, comment)
	bot.Send(tgbotapi.NewMessage(models.TrainerID, text))

	photo := tgbotapi.NewPhoto(models.TrainerID, tgbotapi.FileID(fileID))
	photo.Caption = fmt.Sprintf("📸 Фото от @%s", username)
	bot.Send(photo)

	bot.Send(tgbotapi.NewMessage(userID, "✅ Комментарий и фото отправлены тренеру. Спасибо!"))
}
