package utils

import (
	"os"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitializeLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	Logger.SetFormatter(&logrus.TextFormatter{})
	Logger.SetLevel(logrus.DebugLevel)
}

func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		Logger.Infof("Received %s request for %s", c.Request().Method, c.Request().RequestURI)
		err := next(c)
		Logger.Infof("Responded with status %d", c.Response().Status)
		return err
	}
}
