package handlers

import (
	"fitness-bot/models"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleClientListCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil || update.CallbackQuery.Data != "show_clients" {
		return
	}

	users, err := models.GetAllUsers()
	if err != nil || len(users) == 0 {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "У тебя пока нет зарегистрированных подопечных.")
		bot.Send(msg)
		return
	}

	for _, user := range users {
		text := models.FormatUser(user)
		button := tgbotapi.NewInlineKeyboardButtonData("📂 Карточка", fmt.Sprintf("open_card_%d", user.ID))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(button),
		)

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}
