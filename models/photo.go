package models

import (
	"context"
	"time"
)

type Photo struct {
	ID         int
	UserID     int64
	FileID     string
	Comment    string
	UploadedAt time.Time
}

func SaveUserPhoto(userID int64, fileID string) error {
	_, err := DB.Exec(context.Background(), `
		INSERT INTO photos (user_id, file_id)
		VALUES ($1, $2)
	`, userID, fileID)
	return err
}

func AddPhotoComment(userID int64, comment string) error {
	_, err := DB.Exec(context.Background(), `
		UPDATE photos
		SET comment = $1
		WHERE user_id = $2
		ORDER BY uploaded_at DESC
		LIMIT 1
	`, comment, userID)
	return err
}

func GetUserPhotos(userID int64, limit int) ([]Photo, error) {
	rows, err := DB.Query(context.Background(), `
		SELECT id, file_id, comment, uploaded_at
		FROM photos
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Photo
	for rows.Next() {
		var p Photo
		err := rows.Scan(&p.ID, &p.FileID, &p.Comment, &p.UploadedAt)
		if err != nil {
			continue
		}
		result = append(result, p)
	}

	return result, nil
}
