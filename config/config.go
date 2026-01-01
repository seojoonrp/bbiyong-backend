// config/config.go

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	MongoURI          string
	DBName            string
	JWTSecret         string
	GoogleWebClientID string
	AppleBundleID     string
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	AppConfig = Config{
		Port:              getEnv("PORT", "8080"),
		MongoURI:          getEnv("MONGO_URI", ""),
		DBName:            getEnv("DB_NAME", "bbiyong"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		GoogleWebClientID: getEnv("GOOGLE_WEB_CLIENT_ID", ""),
		AppleBundleID:     getEnv("APPLE_BUNDLE_ID", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
