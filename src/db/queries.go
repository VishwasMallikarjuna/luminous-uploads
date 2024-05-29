package db

import (
	"database/sql"
	"fmt"

	"github.com/VishwasMallikarjuna/luminous-uploads/model"
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

// Query for the most popular image format
func GetMostPopularFormat() (string, error) {
	var mostPopularFormat string
	err := DB.QueryRow("SELECT format FROM image_detail GROUP BY format ORDER BY COUNT(*) DESC LIMIT 1;").Scan(&mostPopularFormat)
	if err != nil && err != sql.ErrNoRows {
		return " ", fmt.Errorf("database query error: %v", err)
	}
	if err == sql.ErrNoRows {
		return " ", nil
	}
	return mostPopularFormat, nil
}

// Query for the top 10 most popular camera models
func GetTopCameraModels() ([]model.CameraModelCount, error) {
	rows, err := DB.Query("SELECT camera_model, COUNT(*) AS count FROM image_detail GROUP BY camera_model ORDER BY count DESC LIMIT 10;")
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var topCameraModels []model.CameraModelCount

	// Iterate over the result set and scan into the CameraModelCount struct
	for rows.Next() {
		var cameraModelCount model.CameraModelCount
		err := rows.Scan(&cameraModelCount.CameraModel, &cameraModelCount.Count)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		topCameraModels = append(topCameraModels, cameraModelCount)
	}

	// Check for errors after iterating through the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return topCameraModels, nil

}

// Query for the image upload frequency per day for the past 30 days
func GetUploadFrequency() ([]model.UploadFrequency, error) {

	rows, err := DB.Query("SELECT DATE(upload_timestamp) AS UploadDate, COUNT(*) AS UploadCount FROM image_detail WHERE upload_timestamp >= NOW() - INTERVAL '30 days' GROUP BY DATE(upload_timestamp) ORDER BY UploadDate DESC;")

	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var topUploadFrequency []model.UploadFrequency

	for rows.Next() {
		var uploadFrequency model.UploadFrequency
		err := rows.Scan(&uploadFrequency.UploadDate, &uploadFrequency.UploadCount)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		topUploadFrequency = append(topUploadFrequency, uploadFrequency)
	}

	// Check for errors after iterating through the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return topUploadFrequency, nil

}
