package main

import (
	"log"

	"goawayauth/internal/config"
	"goawayauth/internal/handlers"
	"goawayauth/internal/repository"
	"goawayauth/internal/service"
	"goawayauth/internal/storage"
)

func main() {
	cfg := config.Load()

	db := storage.NewPostgresDB(cfg)
	defer db.Close()

	repo := repository.NewSecretRepository(db)
	srv := service.NewSecretService(repo)
	h := handlers.NewHandler(srv)

	router := handlers.SetupRouter(h)

	log.Printf("🚀 Сервер стартует на http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
