package http

import (
	"auth-service/internal/infrastructure/http/middleware"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"time"
)

func NewRouter() *echo.Echo {
	e := echo.New()
	rateLimiter := middleware.NewRateLimiter(100, time.Minute, 100)

	e.Use(middleware.Compression) // Response compression
	e.Use(middleware.CORS)
	e.Use(rateLimiter.Middleware())
	e.Use(middleware.SetRequestContextWithTimeout(10 * time.Second))

	logFile, err := os.OpenFile("logs/http.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		e.Logger.Fatal("Failed to open log file:", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			e.Logger.Fatal("Failed to close log file:", err)
		}
	}(logFile)

	customLogger := log.New(logFile, "", 0)

	e.Use(middleware.CustomLogger(middleware.LoggerConfig{
		Format:     "[${time}] ${method} ${uri} - ${status} - ${latency} - ${remote_ip}",
		TimeFormat: "2006-01-02 15:04:05",
		SkipPaths:  []string{"/health"},
		Output:     customLogger,
	}))
	return e
}
