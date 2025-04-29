package handlers

import (
	"context"
	"fitness-bot/models"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleDeleteUserCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	if !strings.HasPrefix(data, "delete_user_") {
		return
	}

	userIDStr := strings.TrimPrefix(data, "delete_user_")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Println("❌ Ошибка парсинга userID при удалении:", err)
		return
	}

	_, err = models.DB.Exec(context.Background(), `DELETE FROM users WHERE telegram_id = $1`, userID)
	if err != nil {
		log.Println("❌ Ошибка при удалении пользователя:", err)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Не удалось удалить пользователя 😢")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("✅ Пользователь ID %d удалён.", userID))
	bot.Send(msg)
}
