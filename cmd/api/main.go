package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"goawayauth/internal/handlers"
	"goawayauth/internal/repository"
	"goawayauth/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvOrFail(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Критическая ошибка: переменная %s не задана!", key)
	}
	return val
}

func loadConfig() *Config {
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

func initDB(cfg *Config) *sql.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBAddr, cfg.DBPort, cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}

	log.Println("Успешное подключение к БД!")
	return db
}

func setupRouter(h *handlers.Handler) *gin.Engine {
	router := gin.Default()

	h.InitRoutes(router)

	return router
}

func main() {
	cfg := loadConfig()

	db := initDB(cfg)
	defer db.Close()

	repo := repository.NewSecretRepository(db)

	service := service.NewSecretService(repo)

	myHandler := handlers.NewHandler(service)

	router := setupRouter(myHandler)

	port := cfg.Port

	log.Printf("🚀 Сервер стартует на http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
