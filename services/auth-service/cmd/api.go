package cmd

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {

		return value

	}
	return fallback
}
