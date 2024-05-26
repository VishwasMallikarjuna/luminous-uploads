package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/VishwasMallikarjuna/luminous-uploads/model"
	"github.com/VishwasMallikarjuna/luminous-uploads/utils"
)

func GenerateUploadLink(c echo.Context) error {
	duration := c.Param("duration")
	expiration, err := time.ParseDuration(duration)
	if err != nil {
		utils.Logger.Errorf("Invalid expiration time: %v", err)
		return c.String(http.StatusBadRequest, "Invalid expiration time")
	}

	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	if authHeader == "" {
		utils.Logger.Warn("Authorization header missing")
		return c.String(http.StatusUnauthorized, "Authorization header missing")
	}

	_, err = ValidateToken(authHeader)
	if err != nil {
		utils.Logger.Errorf("Token validation failed: %v", err)
		return c.String(http.StatusUnauthorized, err.Error())
	}

	uploadLink := generateExpirableLink(expiration)
	utils.Logger.Infof("Generated upload link")
	return c.JSON(http.StatusOK, uploadLink)
}

func generateExpirableLink(expiration time.Duration) model.UploadLink {
	expirationTime := time.Now().Add(expiration)
	url := "/upload-image"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})
	tokenString, err := token.SignedString([]byte(utils.AppConfig.SecretKey))
	if err != nil {
		utils.Logger.Fatalf("Failed to sign token: %v", err)
	}

	return model.UploadLink{
		URL:       url,
		ExpiresAt: expirationTime,
		Token:     tokenString,
	}
}
