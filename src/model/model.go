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
