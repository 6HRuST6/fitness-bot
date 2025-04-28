package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env, если он есть
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Файл .env не найден, продолжаем без него...")
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatal("❌ Нет TELEGRAM_TOKEN в окружении")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("❌ Нет DATABASE_URL в окружении")
	}

	// Дальше идёт подключение базы и запуск бота...
}
