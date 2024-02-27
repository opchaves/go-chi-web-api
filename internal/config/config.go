package config

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var doOnce sync.Once

var (
	Name    = getEnv("APP_NAME", "kommonei")
	Env     = getEnv("APP_ENV", "development")
	Host    = getEnv("HOST", "0.0.0.0")
	Port    = getEnv("PORT", "8080")
	Origins = getEnv("ORIGINS", "")

	IsProduction = Env == "production"
	IsLocal      = Env == "development" || Env == "test"
	DatabaseURL  = getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/app_dev?sslmode=disable")

	JwtSecret        = getEnv("JWT_SECRET", "superSecret")
	JwtExpiry        = toDuration("JWT_EXPIRY", "15m")
	JwtRefreshExpiry = toDuration("JWT_REFRESH_EXPIRY", "1h")

	LoginTokenURL    = getEnv("LOGIN_TOKEN_URL", "http://localhost:8080/auth/token")
	LoginTokenLength = toInt("LOGIN_TOKEN_LENGTH", "16")
	LoginTokenExpiry = toDuration("LOGIN_TOKEN_EXPIRY", "10m")

	EmailSmtpHost = ""
	EmailSmtpPort = 0
	EmailSmtpUser = ""
	EmailSmptPass = ""
	EmailFromAddr = ""
	EmailFromName = ""

	LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func toDuration(envVar string, defaultVal string) time.Duration {
	val, err := time.ParseDuration(getEnv(envVar, defaultVal))
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", envVar, err)
	}
	return val
}

func toInt(envVar string, defaultVal string) int {
	val, err := strconv.Atoi(getEnv(envVar, defaultVal))
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", envVar, err)
	}
	return val
}

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
