package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var doOnce sync.Once

var (
	Name = getEnv("APP_NAME", "kommonei")
	Env  = getEnv("APP_ENV", "development")
	Host = getEnv("HOST", "0.0.0.0")
	Port = getEnv("PORT", "8080")

	IsProduction = Env == "production"
	IsLocal      = Env == "development" || Env == "test"
	DatabaseURL  = getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/app_dev?sslmode=disable")
)

func getEnv(name, defaultValue string) string {
	doOnce.Do(func() {
		readEnvFile(".env")
	})

	if value := os.Getenv(name); value != "" {
		return value
	}

	return defaultValue
}

func readEnvFile(filename string) {
	env := os.Getenv("APP_ENV")
	if env != "production" {
		err := godotenv.Load(filename)
		if err != nil {
			log.Printf("No %s file found. Using default values.\n", filename)
		}
	}
}
