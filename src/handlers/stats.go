package handlers

import (
	"net/http"

	"github.com/VishwasMallikarjuna/luminous-uploads/db"
	"github.com/VishwasMallikarjuna/luminous-uploads/model"
	"github.com/VishwasMallikarjuna/luminous-uploads/utils"
	"github.com/labstack/echo"
)

func GetImageStats(c echo.Context) error {

	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	if authHeader == "" {
		utils.Logger.Warn("Authorization header missing")
		return c.String(http.StatusUnauthorized, "Authorization header missing")
	}

	_, err := ValidateToken(authHeader)
	if err != nil {
		utils.Logger.Errorf("Token validation failed: %v", err)
		return c.String(http.StatusUnauthorized, err.Error())
	}

	mostPopularFormat, err := db.GetMostPopularFormat()
	if err != nil {
		utils.Logger.Errorf("Failed to retrieve mostPopularFormat from the database: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	topCameraModels, err := db.GetTopCameraModels()
	if err != nil {
		utils.Logger.Errorf("Failed to retrieve topCameraModels  from the database: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	uploadFrequency, err := db.GetUploadFrequency()
	if err != nil {
		utils.Logger.Errorf("Failed to retrieve uploadFrequency from the database: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	stats := model.ServiceStatistics{
		MostPopularFormat: mostPopularFormat,
		TopCameraModels:   topCameraModels,
		UploadFrequency:   uploadFrequency,
	}

	utils.Logger.Infof("Successfully retrieved stats")

	return c.JSON(http.StatusOK, stats)
}
