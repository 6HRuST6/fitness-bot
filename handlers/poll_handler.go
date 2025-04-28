package handlers

import (
	"fitness-bot/models"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PollState struct {
	Stage    string // "question" → "options"
	Question string
	Options  []string
}

var activePollData *PollState
var activePoll = make(map[int64]*PollState)          // trainerID -> state
var pollVotes = make(map[string]map[string]int)      // question hash → option → count
var votedUsers = make(map[string]map[int64]bool)     // qHash → userID → true
var votedDetails = make(map[string]map[int64]string) // qHash → userID → option

func HandlePollStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	from := update.CallbackQuery.From
	chat := update.CallbackQuery.Message.Chat

	startPoll(bot, chat.ID, from.ID)
}

// вызывается, если пользователь нажал кнопку "📊 Опрос" (в reply-клавиатуре)
func StartPollFromText(bot *tgbotapi.BotAPI, chatID int64, trainerID int64) {
	startPoll(bot, chatID, trainerID)
}

// общая логика запуска
func startPoll(bot *tgbotapi.BotAPI, chatID, trainerID int64) {
	activePoll[trainerID] = &PollState{Stage: "question"}

	msg := tgbotapi.NewMessage(chatID, "✍️ Введите текст вопроса для опроса:")
	bot.Send(msg)
}

func CheckAndHandlePoll(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	trainerID := update.Message.Chat.ID
	state, ok := activePoll[trainerID]
	if !ok {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	switch state.Stage {
	case "question":
		state.Question = text
		state.Stage = "options"
		bot.Send(tgbotapi.NewMessage(trainerID, "✅ Вопрос сохранён. Теперь отправьте варианты ответа по одному в сообщении.\nНапишите `/done`, когда закончите."))
		return true

	case "options":
		if text == "/done" {
			if len(state.Options) < 2 {
				bot.Send(tgbotapi.NewMessage(trainerID, "❗ Нужно минимум 2 варианта ответа."))
				return true
			}
			sendPoll(bot, trainerID, state)
			delete(activePoll, trainerID)
			return true
		}

		state.Options = append(state.Options, text)
		bot.Send(tgbotapi.NewMessage(trainerID, fmt.Sprintf("Вариант добавлен: %s", text)))
		return true
	}

	return false
}

func sendPoll(bot *tgbotapi.BotAPI, trainerID int64, poll *PollState) {

	activePollData = poll

	users, err := models.GetAllUsers()
	if err != nil {
		log.Println("❌ Не удалось получить пользователей:", err)
		bot.Send(tgbotapi.NewMessage(trainerID, "Ошибка при рассылке опроса."))
		return
	}

	for _, user := range users {
		if user.ID == models.TrainerID {
			continue
		}

		buttons := make([]tgbotapi.InlineKeyboardButton, 0, len(poll.Options))
		for i, option := range poll.Options {
			callback := fmt.Sprintf("vote_%s_%d", hash(poll.Question), i)
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(option, callback))
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))

		msg := tgbotapi.NewMessage(user.ID, fmt.Sprintf("📊 %s", poll.Question))
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}

	bot.Send(tgbotapi.NewMessage(trainerID, "✅ Опрос отправлен всем подопечным."))
	pollVotes[hash(poll.Question)] = make(map[string]int)
}
func HandlePollVote(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	if !strings.HasPrefix(data, "vote_") {
		return
	}

	parts := strings.SplitN(data, "_", 3)
	if len(parts) != 3 {
		return
	}

	qHash := parts[1]
	optionIndexStr := parts[2]

	if activePollData == nil || hash(activePollData.Question) != qHash {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "⛔ Опрос устарел или не найден."))
		return
	}

	optionIndex, err := strconv.Atoi(optionIndexStr)
	if err != nil || optionIndex < 0 || optionIndex >= len(activePollData.Options) {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Неверный вариант ответа."))
		return
	}

	userID := update.CallbackQuery.From.ID

	if votedUsers[qHash] == nil {
		votedUsers[qHash] = make(map[int64]bool)
	}

	if votedUsers[qHash][userID] {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "⚠️ Вы уже голосовали в этом опросе."))
		return
	}

	votedUsers[qHash][userID] = true
	if votedDetails[qHash] == nil {
		votedDetails[qHash] = make(map[int64]string)
	}
	option := activePollData.Options[optionIndex]

	votedDetails[qHash][userID] = option

	if pollVotes[qHash] == nil {
		pollVotes[qHash] = make(map[string]int)
	}
	pollVotes[qHash][option]++

	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "✅ Спасибо за голос!"))

	// Удалить кнопки (чтобы не голосовал снова)
	bot.Request(tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{},
	))
}
func hash(s string) string {
	if len(s) == 0 {
		return "0"
	}
	return fmt.Sprintf("%x", len(s)+int(s[0]))
}
func HandleResultsCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Chat.ID != models.TrainerID {
		return
	}

	if activePollData == nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Нет активного опроса."))
		return
	}

	q := activePollData.Question
	opts := activePollData.Options
	qHash := hash(q)

	votes := pollVotes[qHash]
	if votes == nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Голоса ещё не собраны."))
		return
	}

	// Подсчёт
	total := 0
	for _, count := range votes {
		total += count
	}

	if total == 0 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⛔ Опрос есть, но голосов пока нет."))
		return
	}

	text := fmt.Sprintf("📊 *Результаты опроса:*\n\n❓ *%s*\n\n", q)

	for _, opt := range opts {
		count := votes[opt]
		percent := float64(count) / float64(total) * 100
		text += fmt.Sprintf("• %s — %d голосов (%.1f%%)\n", opt, count, percent)
	}

	text += fmt.Sprintf("\n🧮 Всего голосов: %d", total)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func HandleVotersCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Chat.ID != models.TrainerID {
		return
	}

	if activePollData == nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⛔ Нет активного опроса."))
		return
	}

	q := activePollData.Question
	qHash := hash(q)

	details := votedDetails[qHash]
	if details == nil || len(details) == 0 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❗ Ещё никто не голосовал."))
		return
	}

	text := fmt.Sprintf("📋 *Кто голосовал в опросе:*\n\n❓ *%s*\n\n", q)

	for userID, option := range details {
		user := models.GetUser(userID)
		name := fmt.Sprintf("ID %d", userID)
		if user != nil && user.Username != "" {
			name = fmt.Sprintf("@%s", user.Username)
		}
		text += fmt.Sprintf("• %s — %s\n", name, option)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
