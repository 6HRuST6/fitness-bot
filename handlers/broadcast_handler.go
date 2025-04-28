package handlers

import (
	"fitness-bot/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var pendingBroadcast = make(map[int64]bool)

func HandleBroadcastStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID
	startBroadcast(bot, chatID, userID)
}

// для текстовой команды
func StartBroadcastFromText(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
	startBroadcast(bot, chatID, userID)
}

// универсальная логика
func startBroadcast(bot *tgbotapi.BotAPI, chatID int64, trainerID int64) {
	pendingBroadcast[trainerID] = true

	msg := tgbotapi.NewMessage(chatID, "✍️ Введите текст объявления, которое будет отправлено всем подопечным.")
	bot.Send(msg)
}

func CheckAndHandleBroadcast(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	fromID := update.Message.Chat.ID
	if !pendingBroadcast[fromID] {
		return false
	}
	delete(pendingBroadcast, fromID)

	text := fmt.Sprintf("📢 Объявление от тренера:\n\n%s", update.Message.Text)

	users, err := models.GetAllUsers()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(fromID, "❌ Ошибка получения пользователей."))
		return true
	}

	for _, user := range users {
		if user.ID != models.TrainerID {
			msg := tgbotapi.NewMessage(user.ID, text)
			bot.Send(msg)
		}
	}

	bot.Send(tgbotapi.NewMessage(fromID, "✅ Объявление разослано всем подопечным."))
	return true
}
