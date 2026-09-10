package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func getEnvOrFail(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Критическая ошибка: переменная %s не задана!", key)
	}
	return val
}

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Файл .env не найден, используем системные переменные")
	}

	dbUser := getEnvOrFail("DB_USER")
	dbPass := getEnvOrFail("DB_PASS")
	dbAddr := getEnvOrFail("DB_ADDRESS")
	dbPort := getEnvOrFail("DB_PORT")
	dbName := getEnvOrFail("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbAddr, dbPort, dbName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}

	log.Println("Успешное подключение к БД GoAwayAuth!")

}
