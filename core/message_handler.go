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
	text := update.Message.Text

	// --- Обработка кнопок клиента ---
	switch text {
	case "📸 Отправить фото":
		pendingPhotoRequest[chatID] = true
		Bot.Send(tgbotapi.NewMessage(chatID, "📷 Пожалуйста, пришли фото."))
		return

	case "✍️ Добавить комментарий":
		pendingCommentRequest[chatID] = true
		Bot.Send(tgbotapi.NewMessage(chatID, "✍️ Напиши комментарий, и я передам его тренеру."))
		return
	}

	// --- Обработка комментария ---
	if pendingCommentRequest[chatID] {
		delete(pendingCommentRequest, chatID)

		user := update.Message.From
		comment := update.Message.Text

		username := user.UserName
		if username == "" {
			username = "без username"
		}

		commentMsg := fmt.Sprintf("✍️ Комментарий от @%s (%s):\n\n%s", username, user.FirstName, comment)
		Bot.Send(tgbotapi.NewMessage(models.TrainerID, commentMsg))
		Bot.Send(tgbotapi.NewMessage(chatID, "✅ Комментарий отправлен тренеру!"))
		return
	}

	// --- Обработка фото ---
	if update.Message.Photo != nil {
		if pendingPhotoRequest[chatID] {
			delete(pendingPhotoRequest, chatID)

			fileID := update.Message.Photo[len(update.Message.Photo)-1].FileID
			err := models.SaveUserPhoto(chatID, fileID)
			if err != nil {
				log.Println("Ошибка сохранения фото:", err)
				Bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка при сохранении фото."))
			} else {
				Bot.Send(tgbotapi.NewMessage(chatID, "✅ Фото сохранено!"))
			}
			return
		}

		Bot.Send(tgbotapi.NewMessage(chatID, "ℹ️ Сначала нажми 📸 Отправить фото."))
		return
	}

	// --- Остальные встроенные функции ---
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

	// --- Команды: старт, меню, клиенты ---
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
			msg := tgbotapi.NewMessage(userID, "Ты зарегистрирован ✅ Выбери действие:")
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
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(button),
			)

			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyMarkup = keyboard
			Bot.Send(msg)
		}
		return
	}

	// --- Нераспознанное сообщение ---
	msg := tgbotapi.NewMessage(chatID, "Я пока не понимаю это сообщение 😅")
	if chatID != models.TrainerID {
		msg.ReplyMarkup = clientKeyboard()
	}
	Bot.Send(msg)
}
