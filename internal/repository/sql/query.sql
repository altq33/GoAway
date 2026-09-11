-- name: CreateSecret :exec
INSERT INTO secrets (
    id, encrypted_text, is_client_encrypted, password_hash, notify_email, file_path, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetSecretMeta :one
SELECT 
    id, 
    (password_hash IS NOT NULL)::boolean AS has_password, 
    is_client_encrypted
FROM secrets 
WHERE id = $1 LIMIT 1;

-- name: GetSecretForRead :one
SELECT * FROM secrets 
WHERE id = $1 LIMIT 1;

-- name: DeleteSecret :exec
-- Удаляем записку после успешного прочтения.
DELETE FROM secrets WHERE id = $1;

-- name: GetExpiredSecrets :many
SELECT id, file_path FROM secrets 
WHERE expires_at < NOW();

-- name: CleanExpiredSecrets :exec
DELETE FROM secrets 
WHERE expires_at < NOW();