package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/VishwasMallikarjuna/luminous-uploads/db"
	"github.com/VishwasMallikarjuna/luminous-uploads/utils"
)

func UploadImage(c echo.Context) error {
	// Validate the token
	_, err := ValidateUploadImageToken(c)
	if err != nil {
		return err
	}

	// Process the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		utils.Logger.Errorf("No file is received: %v", err)
		return c.String(http.StatusBadRequest, "No file is received")
	}

	files := form.File["images"]
	var ids []int

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			utils.Logger.Errorf("Unable to open the file: %v", err)
			return c.String(http.StatusInternalServerError, "Unable to open the file")
		}
		defer src.Close()

		imageData, err := io.ReadAll(src)
		if err != nil {
			utils.Logger.Errorf("Unable to read the file: %v", err)
			return c.String(http.StatusInternalServerError, "Unable to read the file")
		}

		imageHash := utils.CalculateHash(imageData)
		existingID, err := db.CheckImageHash(imageHash)
		if err != nil {
			utils.Logger.Errorf("Failed to check image hash: %v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if existingID != 0 {
			utils.Logger.Infof("Image already exists with ID: %d", existingID)
			ids = append(ids, existingID)
			continue
		}

		insertedID, err := db.InsertImage(imageData, imageHash)
		if err != nil {
			utils.Logger.Errorf("Failed to store the image in the database: %v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		_, err = src.Seek(0, io.SeekStart)
		if err != nil {
			utils.Logger.Errorf("Failed to reset file pointer: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to reset file pointer")
		}

		insertImageMetadata, err := utils.ExtractMetadata(src)
		if err != nil {
			utils.Logger.Errorf("Failed to extract image metadata: %v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		err = db.InsertImageMetadata(insertedID, insertImageMetadata.Width, insertImageMetadata.Height, insertImageMetadata.CameraModel, insertImageMetadata.Location, insertImageMetadata.Format)
		if err != nil {
			utils.Logger.Errorf("Failed to store the image metadata in the database: %v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		utils.Logger.Infof("Successfully uploaded image with ID: %d", insertedID)
		ids = append(ids, insertedID)
	}

	return c.JSON(http.StatusOK, ids)
}

func GetImage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		utils.Logger.Errorf("Invalid image ID: %v", err)
		return c.String(http.StatusBadRequest, "Invalid image ID")
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

	imageData, err := db.GetImageByID(id)
	if err != nil {
		if err.Error() == "image not found" {
			utils.Logger.Warn("Image not found")
			return c.String(http.StatusNotFound, "Image not found")
		}
		utils.Logger.Errorf("Failed to retrieve image from the database: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	utils.Logger.Infof("Successfully retrieved image with ID: %d", id)
	return c.Blob(http.StatusOK, "image/jpeg", imageData)
}
