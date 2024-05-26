package main

import (
	"flag"
	"os"

	"github.com/labstack/echo"
	"github.com/VishwasMallikarjuna/luminous-uploads/db"
	"github.com/VishwasMallikarjuna/luminous-uploads/handlers"
	"github.com/VishwasMallikarjuna/luminous-uploads/utils"
)

func main() {
	// Define a flag for the configuration file path
	utils.LoadConfig()
	flag.Parse() // Parse the flags

	utils.InitializeLogger()
	// Load configuration using the provided file path
	utils.LoadConfig() // Adjusted to pass the path as an argument
	
	db.Connect() // Connect to the database

	e := echo.New()

	// Middleware to log each request
	e.Use(utils.LogRequest)

	// Define routes
	e.POST("/generate-upload-link/:duration", handlers.GenerateUploadLink)
	e.POST("/upload-image", handlers.UploadImage)
	e.GET("/image/:imageId", handlers.GetImage)

	// Start the server
	err := e.Start(":1323")
	if err != nil {
		e.Logger.Fatal(err)
		os.Exit(2)
	}
}
