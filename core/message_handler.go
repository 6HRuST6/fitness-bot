package core

import (
	"fitness-bot/handlers"
	"fitness-bot/models"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMessage(update tgbotapi.Update) {
	if handlers.CheckAndHandlePoll(Bot, update) {
		return
	}
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// ✍️ Комментарий к фото
	if fileID, ok := handlers.PendingComment(chatID); ok {
		handlers.SavePhotoComment(Bot, update, fileID)
		return
	}

	// ✍️ Проверка: если это рекомендация
	if handlers.CheckAndHandleRecommendation(Bot, update) {
		return
	}

	// 📷 Фото
	if update.Message.Photo != nil {
		userID := update.Message.From.ID
		fileID := update.Message.Photo[len(update.Message.Photo)-1].FileID

		err := models.SaveUserPhoto(userID, fileID)
		if err != nil {
			log.Println("Ошибка сохранения фото:", err)
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при сохранении фото.")
			Bot.Send(msg)
			return
		}

		handlers.MarkPendingPhoto(chatID, fileID)

		msg := tgbotapi.NewMessage(chatID, "✍️ Напиши комментарий к фото. Например: завтрак, обед, самочувствие и т.д.")
		Bot.Send(msg)
		return
	}

	// 💬 Команды/встроенные кнопки
	if handlers.CheckAndHandleBroadcast(Bot, update) {
		return
	}
	if text == "/results" {
		handlers.HandleResultsCommand(Bot, update)
		return
	}
	if text == "/voters" {
		handlers.HandleVotersCommand(Bot, update)
		return
	}
	if text == "📣 Объявление" {
		handlers.StartBroadcastFromText(Bot, chatID, update.Message.From.ID)
		return
	}
	if text == "📊 Опрос" {
		handlers.StartPollFromText(Bot, chatID, update.Message.From.ID)
		return
	}

	switch text {
	case "/start":
		userID := chatID
		username := update.Message.From.UserName
		name := update.Message.From.FirstName

		models.RegisterUser(userID, username, name)

		if userID == models.TrainerID {
			msg := tgbotapi.NewMessage(userID, "📋 Главное меню тренера:")
			msg.ReplyMarkup = trainerKeyboard()
			Bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(userID, "Ты зарегистрирован ✅ Отправляй фото еды или отчёты!")
			Bot.Send(msg)
		}
		return

	case "/menu":
		if chatID == models.TrainerID {
			msg := tgbotapi.NewMessage(chatID, "📋 Главное меню тренера:")
			msg.ReplyMarkup = trainerKeyboard()
			Bot.Send(msg)
		}
		return

	case "/clients":
		if chatID != models.TrainerID {
			msg := tgbotapi.NewMessage(chatID, "⛔ Только тренер может просматривать список подопечных.")
			Bot.Send(msg)
			return
		}

		users, err := models.GetAllUsers()
		if err != nil || len(users) == 0 {
			msg := tgbotapi.NewMessage(chatID, "У тебя пока нет зарегистрированных подопечных.")
			Bot.Send(msg)
			return
		}

		for _, user := range users {
			text := models.FormatUser(user)
			button := tgbotapi.NewInlineKeyboardButtonData("📂 Карточка", fmt.Sprintf("open_card_%d", user.ID))
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(button),
			)

			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyMarkup = keyboard
			Bot.Send(msg)
		}
		return
	}

	// 🤷 Всё остальное
	msg := tgbotapi.NewMessage(chatID, "Я пока не понимаю это сообщение 😅")
	Bot.Send(msg)
}
