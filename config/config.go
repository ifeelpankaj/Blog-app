package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env           string
	Port          string
	DB_user       string
	DB_port       string
	DB_host       string
	DB_pass       string
	DB_name       string
	AllowedOrigin string
}

var AppConfig *Config

func LoadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // default
	}

	envFile := ".env." + env
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Error loading %s file: %v\n", envFile, err)
	}

	AppConfig = &Config{
		Port:          getEnv("PORT", "4000"),
		Env:           getEnv("APP_ENV", "dev"),
		DB_user:       getEnv("DB_USER", "postgres"),
		DB_port:       getEnv("DB_PORT", "5432"),
		DB_host:       getEnv("DB_HOST", "localhost"),
		DB_pass:       getEnv("DB_PASS", "password"),
		DB_name:       getEnv("DB_NAME", "blog_app_db"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
