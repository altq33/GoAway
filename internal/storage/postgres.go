package storage

import (
	"database/sql"
	"fmt"
	"log"

	"goawayauth/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(cfg *config.Config) *sql.DB {
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
