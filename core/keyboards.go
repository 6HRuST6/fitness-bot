package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func trainerKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/clients"),
			tgbotapi.NewKeyboardButton("/results"),
			tgbotapi.NewKeyboardButton("/voters"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📣 Объявление"),
			tgbotapi.NewKeyboardButton("📊 Опрос"),
			tgbotapi.NewKeyboardButton("/menu"),
		),
	)
}
func clientKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📸 Отправить фото"),
			tgbotapi.NewKeyboardButton("✍️ Добавить комментарий"),
		),
	)
}
