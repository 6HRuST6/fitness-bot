package handlers

import (
	"fitness-bot/models"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCardCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	message := update.CallbackQuery.Message

	if !strings.HasPrefix(data, "open_card_") {
		return
	}

	userIDStr := strings.TrimPrefix(data, "open_card_")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Println("❌ Ошибка парсинга userID:", err)
		return
	}

	user := models.GetUser(userID)
	if user == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Пользователь не найден")
		bot.Send(msg)
		return
	}

	text := fmt.Sprintf("Карточка клиента:\n👤 %s (@%s)\nID: %d", user.Name, user.Username, user.ID)

	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🖼 Фото (3)", fmt.Sprintf("view_photos_3_%d", user.ID)),
			tgbotapi.NewInlineKeyboardButtonData("🖼 Фото (5)", fmt.Sprintf("view_photos_5_%d", user.ID)),
			tgbotapi.NewInlineKeyboardButtonData("🖼 Фото (10)", fmt.Sprintf("view_photos_10_%d", user.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📄 Рекомендация", fmt.Sprintf("recommend_%d", user.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Удалить пользователя", fmt.Sprintf("delete_user_%d", user.ID)),
		),
	)

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		message.Chat.ID,
		message.MessageID,
		text,
		buttons,
	)

	_, err = bot.Send(edit)
	if err != nil {
		log.Println("❌ Ошибка при редактировании карточки:", err)
	}

}
