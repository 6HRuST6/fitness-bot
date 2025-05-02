package core

import (
	"fitness-bot/handlers"
	"log"
	"strconv"
	"strings"
	

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallback(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data

	switch {
	case strings.HasPrefix(data, "open_card_"):
		handlers.HandleCardCallback(Bot, update)

	case strings.HasPrefix(data, "recommend_"):
		handlers.HandleRecommendCallback(Bot, update)

	case strings.HasPrefix(data, "view_photos_"):
		handlers.HandlePhotoViewCallback(Bot, update)

	case data == "show_clients":
		handlers.HandleClientListCallback(Bot, update)

	case strings.HasPrefix(data, "delete_user_"):
		handlers.HandleDeleteUserCallback(Bot, update)

	case strings.HasPrefix(data, "broadcast_start"):
		log.Println("Received broadcast_start callback")
		handlers.HandleBroadcastStart(Bot, update)

	case strings.HasPrefix(data, "vote_"):
		handlers.HandlePollVote(Bot, update)

	case data == "poll_start":
		handlers.HandlePollStart(Bot, update)

	case strings.HasPrefix(data, "react_"):
		parts := strings.Split(data, "_")
		if len(parts) != 3 {
			return
		}

		userID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return
		}

		reaction := parts[2]
		var text string

		switch reaction {
		case "good":
			text = "👍 Тренер одобрил фото!"
		case "fire":
			text = "🔥 Отлично! Держи темп!"
		case "warn":
			text = "⚠️ Обрати внимание на питание."
		default:
			text = "👀 Неизвестная реакция."
		}

		// Отправка клиенту
		msg := tgbotapi.NewMessage(userID, text)
		Bot.Send(msg)

		// Подтверждение тренеру
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "✅ Реакция отправлена!")
		Bot.Request(callback)
	}
}
