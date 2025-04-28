package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Временное хранилище: от кого -> кому
var pendingRecommendation = make(map[int64]int64)

func HandleRecommendCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	if !strings.HasPrefix(data, "recommend_") {
		return
	}

	userIDStr := strings.TrimPrefix(data, "recommend_")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Println("❌ Ошибка парсинга userID в recommend:", err)
		return
	}

	fromID := update.CallbackQuery.From.ID
	pendingRecommendation[fromID] = userID

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "✍️ Напиши текст рекомендации:")
	bot.Send(msg)
}

func CheckAndHandleRecommendation(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	fromID := update.Message.Chat.ID
	targetID, ok := pendingRecommendation[fromID]
	if !ok {
		return false
	}

	delete(pendingRecommendation, fromID)

	text := fmt.Sprintf("📩 Рекомендация от тренера:\n\n%s", update.Message.Text)
	msgToClient := tgbotapi.NewMessage(targetID, text)
	_, err := bot.Send(msgToClient)
	if err != nil {
		msg := tgbotapi.NewMessage(fromID, "❌ Не удалось отправить рекомендацию клиенту.")
		bot.Send(msg)
		return true
	}

	msg := tgbotapi.NewMessage(fromID, "✅ Рекомендация отправлена!")
	bot.Send(msg)
	return true
}
