package models

import (
	"context"
	"fmt"
	"log"
	"time"
)

//const TrainerID int64 = 823298509,247753697
var TrainerID = []int64{823298509, 247753697}
type User struct {
	ID       int64
	Username string
	Name     string
	JoinedAt time.Time
}

// ✅ Регистрирует пользователя (если он ещё не существует)
func RegisterUser(id int64, username, name string) {
	_, err := DB.Exec(context.Background(), `
		INSERT INTO users (telegram_id, username, name, joined_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (telegram_id) DO NOTHING
	`, id, username, name, time.Now())
	if err != nil {
		log.Println("❌ Ошибка при регистрации пользователя:", err)
	}
}

// ✅ Получение одного пользователя по ID
func GetUser(id int64) *User {
	row := DB.QueryRow(context.Background(), `
		SELECT telegram_id, username, name, joined_at FROM users WHERE telegram_id = $1
	`, id)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.JoinedAt)
	if err != nil {
		return nil
	}
	return &u
}

// ✅ Форматирование карточки пользователя
func FormatUser(u *User) string {
	return fmt.Sprintf("%s (@%s) — ID: %d — %s",
		u.Name, u.Username, u.ID, u.JoinedAt.Format("02.01.2006 15:04"))
}

// ✅ Получение всех пользователей (для /clients, show_clients )
func GetAllUsers() ([]*User, error) {
	rows, err := DB.Query(context.Background(), `
		SELECT telegram_id, username, name, joined_at
		FROM users
		ORDER BY joined_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Name, &u.JoinedAt); err != nil {
			continue
		}
		result = append(result, &u)
	}
	return result, nil
}
