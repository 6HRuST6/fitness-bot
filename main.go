package main

import (
	"log"
	"os"

	"fitness-bot/core"
	"fitness-bot/models"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Ошибка загрузки .env файла")
	}
	models.InitDB()
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("❌ TELEGRAM_TOKEN не найден в .env")
	}

	core.ClearOldUpdates(token)
	core.Start(token)

	models.InitDB()

}
