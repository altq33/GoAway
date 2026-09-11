package service

import (
	"context"
	"errors"

	"goawayauth/internal/repository"
)

type SecretService interface {
	CreateSecret(ctx context.Context) error
	GetSecretMeta(ctx context.Context, id string) error
	ReadSecret(ctx context.Context, id string, password *string) error
}

type secretService struct {
	repo repository.SecretRepository
}

func NewSecretService(repo repository.SecretRepository) SecretService {
	return &secretService{
		repo: repo,
	}
}

func (s *secretService) CreateSecret(ctx context.Context) error {
	return errors.New("CreateSecret: логика еще не написана")
}

func (s *secretService) GetSecretMeta(ctx context.Context, id string) error {
	return errors.New("GetSecretMeta: логика еще не написана")
}

func (s *secretService) ReadSecret(ctx context.Context, id string, password *string) error {
	return errors.New("ReadSecret: логика еще не написана")
}
