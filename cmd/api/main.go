package main

import (
	"context"
	"log"
	"time"

	"goawayauth/internal/config"
	"goawayauth/internal/handlers"
	"goawayauth/internal/repository"
	"goawayauth/internal/service"
	"goawayauth/internal/storage"
	"goawayauth/internal/worker"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db := storage.NewPostgresDB(cfg)
	defer db.Close()

	repo := repository.NewSecretRepository(db)
	srv := service.NewSecretService(repo)
	h := handlers.NewHandler(srv)

	router := handlers.SetupRouter(h)

	cleaner := worker.NewSecretCleaner(repo)

	go cleaner.Start(ctx, 1*time.Hour)

	log.Printf("🚀 Сервер стартует на http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
