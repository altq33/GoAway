package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := storage.NewPostgresDB(cfg)
	defer db.Close()

	repo := repository.NewSecretRepository(db)
	srv := service.NewSecretService(repo)
	h := handlers.NewHandler(srv)

	r := handlers.SetupRouter(h)

	cleaner := worker.NewSecretCleaner(repo)

	go cleaner.Start(ctx, 1*time.Hour)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Сервер стартует на http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Ошибка запуска сервера", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Получен сигнал завершения. Останавливаем сервер...")

	cancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		slog.Error("Принудительная остановка сервера", "error", err)
	}

	slog.Info("Сервер успешно остановлен. До свидания!")

}
