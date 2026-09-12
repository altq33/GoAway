package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"goawayauth/internal/cryptoutil"
	"goawayauth/internal/repository"
	"goawayauth/internal/repository/db"
)

type SecretService interface {
	CreateSecret(ctx context.Context, req CreateSecretDTO) (CreateSecretResult, error)
	GetSecretMeta(ctx context.Context, id string) error
	ReadSecret(ctx context.Context, id string, password *string) error
}

type secretService struct {
	repo repository.SecretRepository
}

type CreateSecretDTO struct {
	Text              string
	IsClientEncrypted bool
	Password          *string
	TTLHours          int
}

type CreateSecretResult struct {
	ID            string
	EncryptionKey *string
}

func NewSecretService(repo repository.SecretRepository) SecretService {
	return &secretService{
		repo: repo,
	}
}

func (s *secretService) CreateSecret(ctx context.Context, req CreateSecretDTO) (CreateSecretResult, error) {
	// 1. Генерируем ID (16 байт = 32 символа в hex)
	id, err := cryptoutil.GenerateHexID(16)
	if err != nil {
		return CreateSecretResult{}, err
	}

	// 2. Хэшируем пароль (если передан)
	var passwordHash *string
	if req.Password != nil && *req.Password != "" {
		passwordHash, err = cryptoutil.HashPassword(*req.Password)
		if err != nil {
			return CreateSecretResult{}, err
		}
	}

	// 3. Гибридное шифрование
	var encryptionKey *string
	textToStore := req.Text

	if !req.IsClientEncrypted {
		encryptedText, generatedKey, err := cryptoutil.EncryptTextAES(req.Text)
		if err != nil {
			return CreateSecretResult{}, err
		}
		textToStore = encryptedText
		encryptionKey = &generatedKey
	}

	expiresAt := time.Now().Add(time.Duration(req.TTLHours) * time.Hour)

	err = s.repo.CreateSecret(ctx, db.CreateSecretParams{ // <-- Убедись, что пакет db импортирован правильно
		ID:                id,
		EncryptedText:     textToStore,
		IsClientEncrypted: req.IsClientEncrypted,
		PasswordHash:      toNullString(passwordHash), // Используем наш хелпер
		ExpiresAt:         expiresAt,
	})

	if err != nil {
		return CreateSecretResult{}, err
	}

	// 5. Успешный ответ
	return CreateSecretResult{
		ID:            id,
		EncryptionKey: encryptionKey,
	}, nil
}

func (s *secretService) GetSecretMeta(ctx context.Context, id string) error {
	return errors.New("GetSecretMeta: логика еще не написана")
}

func (s *secretService) ReadSecret(ctx context.Context, id string, password *string) error {
	return errors.New("ReadSecret: логика еще не написана")
}

// маленькая утилита-хелпер для чистой конвертации
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
