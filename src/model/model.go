package model

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type UploadLink struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
}

type CustomClaims struct {
	jwt.StandardClaims
}

type ImageMetadata struct {
	Width       int
	Height      int
	CameraModel string
	Location    string
	Format      string
}

type ServiceStatistics struct {
	MostPopularFormat string             `json:"most_popular_format"`
	TopCameraModels   []CameraModelCount `json:"top_camera_models"`
	UploadFrequency   []UploadFrequency  `json:"upload_frequency"`
}

type CameraModelCount struct {
	CameraModel string `json:"camera_model"`
	Count       int    `json:"count"`
}

type UploadFrequency struct {
	UploadDate  string `json:"upload_date"`
	UploadCount int    `json:"upload_count"`
}
