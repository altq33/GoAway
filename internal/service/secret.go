package service

import (
	"context"
	"errors"
	"time"

	"goawayauth/internal/cryptoutil"
	"goawayauth/internal/repository"
	"goawayauth/internal/repository/db"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type SecretService interface {
	CreateSecret(ctx context.Context, req CreateSecretDTO) (CreateSecretResult, error)
	GetSecretMeta(ctx context.Context, id string) (SecretMetaResult, error)
	ReadSecret(ctx context.Context, req ReadSecretDTO) (string, error)
}

type secretService struct {
	repo repository.SecretRepository
}

type SecretMetaResult struct {
	IsClientEncrypted bool
	RequiresPassword  bool
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

type ReadSecretDTO struct {
	ID            string
	Password      *string
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

	err = s.repo.CreateSecret(ctx, db.CreateSecretParams{
		ID:                id,
		EncryptedText:     textToStore,
		IsClientEncrypted: req.IsClientEncrypted,
		PasswordHash:      db.ToNullString(passwordHash),
		ExpiresAt:         expiresAt,
	})

	if err != nil {
		return CreateSecretResult{}, err
	}

	return CreateSecretResult{
		ID:            id,
		EncryptionKey: encryptionKey,
	}, nil
}

func (s *secretService) checkExpiration(ctx context.Context, id string, expiresAt time.Time) error {
	if time.Now().After(expiresAt) {
		_ = s.repo.DeleteSecret(ctx, id)
		return errors.New("срок действия записки истек")
	}
	return nil
}

func (s *secretService) ReadSecret(ctx context.Context, req ReadSecretDTO) (string, error) {
	secret, err := s.repo.GetSecretForRead(ctx, req.ID)
	if err != nil {
		return "", errors.New("записка не найдена или уже была прочитана")
	}
	if err := s.checkExpiration(ctx, secret.ID, secret.ExpiresAt); err != nil {
		return "", err
	}

	if secret.PasswordHash.Valid { // Вспоминаем нашу sql.NullString
		if req.Password == nil || *req.Password == "" {
			return "", errors.New("требуется пароль")
		}

		err = bcrypt.CompareHashAndPassword([]byte(secret.PasswordHash.String), []byte(*req.Password))
		if err != nil {
			return "", errors.New("неверный пароль")
		}
	}

	var resultText string = secret.EncryptedText

	if !secret.IsClientEncrypted {

		if req.EncryptionKey == nil || *req.EncryptionKey == "" {
			return "", errors.New("для расшифровки требуется ключ (encryptionKey)")
		}

		decryptedText, err := cryptoutil.DecryptTextAES(secret.EncryptedText, *req.EncryptionKey)
		if err != nil {
			return "", errors.New("не удалось расшифровать записку: " + err.Error())
		}
		resultText = decryptedText

	}

	err = s.repo.DeleteSecret(ctx, req.ID)
	if err != nil {
		slog.Error("Не удалось удалить прочитанную записку",
			slog.String("secret_id", req.ID),
			slog.Any("error", err),
		)
	}

	return resultText, nil
}

func (s *secretService) GetSecretMeta(ctx context.Context, id string) (SecretMetaResult, error) {
	secret, err := s.repo.GetSecretMeta(ctx, id)
	if err != nil {
		return SecretMetaResult{}, errors.New("записка не найдена или уже была прочитана")
	}

	if err := s.checkExpiration(ctx, id, secret.ExpiresAt); err != nil {
		return SecretMetaResult{}, err
	}

	return SecretMetaResult{
		IsClientEncrypted: secret.IsClientEncrypted,
		RequiresPassword:  secret.HasPassword,
	}, nil
}
