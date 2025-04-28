package models

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL не задан в .env")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к базе: %v", err)
	}

	log.Println("✅ Подключение к PostgreSQL установлено")
}
