package service

import (
	"context"
	"errors"

	"goawayauth/internal/repository"
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
	return CreateSecretResult{}, errors.New("CreateSecret: логика еще не написана")
}

func (s *secretService) GetSecretMeta(ctx context.Context, id string) error {
	return errors.New("GetSecretMeta: логика еще не написана")
}

func (s *secretService) ReadSecret(ctx context.Context, id string, password *string) error {
	return errors.New("ReadSecret: логика еще не написана")
}
