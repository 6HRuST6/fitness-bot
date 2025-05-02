package core

import (
	"fitness-bot/handlers"
	"fitness-bot/models"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var pendingPhotoRequest = make(map[int64]bool)
var pendingCommentRequest = make(map[int64]bool)

func handleMessage(update tgbotapi.Update) {
	if handlers.CheckAndHandlePoll(Bot, update) {
		return
	}
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	user := update.Message.From
	text := update.Message.Text

	// === Обработка кнопок ===
	switch text {
	case "📸 Отправить фото":
		pendingPhotoRequest[chatID] = true
		Bot.Send(tgbotapi.NewMessage(chatID, "📷 Пришли, пожалуйста, фото."))
		return

	case "✍️ Добавить комментарий":
		pendingCommentRequest[chatID] = true
		Bot.Send(tgbotapi.NewMessage(chatID, "✍️ Напиши комментарий, и я передам его тренеру."))
		return
	}

	// === Обработка фото ===
	if update.Message.Photo != nil {
		fileID := update.Message.Photo[len(update.Message.Photo)-1].FileID

		// Сохраняем в базу (опционально)
		err := models.SaveUserPhoto(chatID, fileID)
		if err != nil {
			log.Println("❌ Ошибка сохранения фото:", err)
			Bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка при сохранении фото."))
		} else {
			username := user.UserName
			if username == "" {
				username = "без username"
			}

			caption := fmt.Sprintf("📷 Фото от @%s (%s)", username, user.FirstName)
			photo := tgbotapi.NewPhoto(models.TrainerID, tgbotapi.FileID(fileID))
			photo.Caption = caption
			Bot.Send(photo)

			Bot.Send(tgbotapi.NewMessage(chatID, "✅ Фото сохранено и отправлено тренеру."))
		}
		delete(pendingPhotoRequest, chatID)
		return
	}

	// === Обработка комментариев по кнопке  ===
	if pendingCommentRequest[chatID] {
	delete(pendingCommentRequest, chatID)

	username := user.UserName
	if username == "" {
		username = "без username"
	}

	message := fmt.Sprintf("✍️ Комментарий от @%s (%s):\n\n%s", username, user.FirstName, text)
	Bot.Send(tgbotapi.NewMessage(models.TrainerID, message))
	Bot.Send(tgbotapi.NewMessage(chatID, "✅ Комментарий отправлен тренеру!"))
	return
}
	// === Встроенные функции ===
	if handlers.CheckAndHandleRecommendation(Bot, update) {
		return
	}
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

	// === Команды: старт, меню, клиенты ===
	switch text {
	case "/start":
		models.RegisterUser(chatID, user.UserName, user.FirstName)

		if chatID == models.TrainerID {
			msg := tgbotapi.NewMessage(chatID, "📋 Главное меню тренера:")
			msg.ReplyMarkup = trainerKeyboard()
			Bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Ты зарегистрирован ✅ Выбери действие:")
			msg.ReplyMarkup = clientKeyboard()
			Bot.Send(msg)
		}
		return

	case "/menu":
		if chatID == models.TrainerID {
			msg := tgbotapi.NewMessage(chatID, "📋 Главное меню тренера:")
			msg.ReplyMarkup = trainerKeyboard()
			Bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "📋 Меню пользователя:")
			msg.ReplyMarkup = clientKeyboard()
			Bot.Send(msg)
		}
		return

	case "/clients":
		if chatID != models.TrainerID {
			Bot.Send(tgbotapi.NewMessage(chatID, "⛔ Только тренер может просматривать список подопечных."))
			return
		}

		users, err := models.GetAllUsers()
		if err != nil || len(users) == 0 {
			Bot.Send(tgbotapi.NewMessage(chatID, "У тебя пока нет зарегистрированных подопечных."))
			return
		}

		for _, user := range users {
			text := models.FormatUser(user)
			button := tgbotapi.NewInlineKeyboardButtonData("📂 Карточка", fmt.Sprintf("open_card_%d", user.ID))
			keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))

			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyMarkup = keyboard
			Bot.Send(msg)
		}
		return
	}

	// === Всё остальное ===
	msg := tgbotapi.NewMessage(chatID, "Я пока не понимаю это сообщение 😅")
	if chatID != models.TrainerID {
		msg.ReplyMarkup = clientKeyboard()
	}
	Bot.Send(msg)
}
