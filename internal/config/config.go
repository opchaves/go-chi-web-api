package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var doOnce sync.Once

// GetEnv returns the value of the environment variable named by the key,
// or a default value (defVal) if the environment variable is not present.
func GetEnv(name, defVal string) string {
	doOnce.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found. Using default values.")
		}
	})

	if value := os.Getenv(name); value != "" {
		return value
	}

	return defVal
}

var (
	Name = GetEnv("NAME", "go-chi-web-api")
	Env  = GetEnv("ENV", "development")
	Host = GetEnv("HOST", "0.0.0.0")
	Port = GetEnv("PORT", "8080")

	IsProduction = Env == "production"
	IsLocal      = Env == "development" || Env == "test"
)
