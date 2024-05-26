package db

import (
	"database/sql"
	"fmt"
)

func InsertImage(imageData []byte, imageHash string) (int, error) {
	var insertedID int
	err := DB.QueryRow("INSERT INTO images (image_data, image_hash) VALUES ($1, $2) RETURNING id", imageData, imageHash).Scan(&insertedID)
	if err != nil {
		return 0, fmt.Errorf("failed to store the image in the database: %v", err)
	}
	return insertedID, nil
}

func CheckImageHash(imageHash string) (int, error) {
	var existingID int
	err := DB.QueryRow("SELECT id FROM images WHERE image_hash = $1", imageHash).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("database query error: %v", err)
	}
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return existingID, nil
}

func InsertImageMetadata(id, width, height int, cameraModel, location, format string) error {
	_, err := DB.Exec("INSERT INTO image_detail (id, width, height, camera_model, location, format) VALUES ($1, $2, $3, $4, $5, $6)", id, width, height, cameraModel, location, format)
	if err != nil {
		return fmt.Errorf("failed to store the image metadata in the database: %v", err)
	}
	return nil
}

func GetImageByID(id int) ([]byte, error) {
	var imageData []byte
	err := DB.QueryRow("SELECT image_data FROM images WHERE id = $1", id).Scan(&imageData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to retrieve the image: %v", err)
	}
	return imageData, nil
}
