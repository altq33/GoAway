package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser string
	DBPass string
	DBAddr string
	DBPort string
	DBName string
	Port   string
}

func Load() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Файл .env не найден, используем системные переменные")
	}

	return &Config{
		DBUser: getEnvOrFail("DB_USER"),
		DBPass: getEnvOrFail("DB_PASS"),
		DBAddr: getEnvOrFail("DB_ADDRESS"),
		DBPort: getEnvOrFail("DB_PORT"),
		DBName: getEnvOrFail("DB_NAME"),
		Port:   getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvOrFail(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Критическая ошибка: переменная %s не задана!", key)
	}
	return val
}
