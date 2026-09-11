package repository

import (
	"context"
	"database/sql"

	"goawayauth/internal/repository/db"
)

type SecretRepository interface {
	CreateSecret(ctx context.Context, arg db.CreateSecretParams) error
	GetSecretMeta(ctx context.Context, id string) (db.GetSecretMetaRow, error)
	GetSecretForRead(ctx context.Context, id string) (db.Secret, error)
	DeleteSecret(ctx context.Context, id string) error
	GetExpiredSecrets(ctx context.Context) ([]db.GetExpiredSecretsRow, error)
	CleanExpiredSecrets(ctx context.Context) error
}

type postgresSecretRepo struct {
	dbConn *sql.DB
	q      *db.Queries
}

func NewSecretRepository(conn *sql.DB) SecretRepository {
	return &postgresSecretRepo{
		dbConn: conn,
		q:      db.New(conn),
	}	
}

func (r *postgresSecretRepo) CreateSecret(ctx context.Context, arg db.CreateSecretParams) error {

	return r.q.CreateSecret(ctx, arg)
}

func (r *postgresSecretRepo) GetSecretMeta(ctx context.Context, id string) (db.GetSecretMetaRow, error) {
	return r.q.GetSecretMeta(ctx, id)
}

func (r *postgresSecretRepo) GetSecretForRead(ctx context.Context, id string) (db.Secret, error) {
	return r.q.GetSecretForRead(ctx, id)
}

func (r *postgresSecretRepo) DeleteSecret(ctx context.Context, id string) error {
	return r.q.DeleteSecret(ctx, id)
}

func (r *postgresSecretRepo) GetExpiredSecrets(ctx context.Context) ([]db.GetExpiredSecretsRow, error) {
	return r.q.GetExpiredSecrets(ctx)
}

func (r *postgresSecretRepo) CleanExpiredSecrets(ctx context.Context) error {
	return r.q.CleanExpiredSecrets(ctx)
}
