package worker

import (
	"context"
	"log/slog"
	"time"

	"goawayauth/internal/repository"
)

type SecretCleaner struct {
	repo repository.SecretRepository
}

func NewSecretCleaner(repo repository.SecretRepository) *SecretCleaner {
	return &SecretCleaner{
		repo: repo,
	}
}

// Start запускает бесконечный цикл очистки в фоне
func (c *SecretCleaner) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info("Воркер очистки запущен", slog.String("интервал", interval.String()))

	for {
		select {
		case <-ctx.Done():
			slog.Info("Воркер очистки остановлен (сигнал завершения)")
			return

		case <-ticker.C:
			err := c.repo.CleanExpiredSecrets(ctx)
			if err != nil {
				slog.Error("Ошибка при фоновой очистке базы", slog.Any("error", err))
			}
		}
	}
}
