package handlers

// ---------------------------------------------------------
// 1. Создание записки (Create)
// ---------------------------------------------------------

type CreateSecretRequest struct {
	Text string `json:"text" binding:"required"`

	IsClientEncrypted bool `json:"isClientEncrypted"`

	Password *string `json:"password,omitempty"`

	TTLHours int `json:"ttlHours" binding:"required,min=1,max=168"`
}

type CreateSecretResponse struct {
	ID string `json:"id"`

	EncryptionKey *string `json:"encryptionKey,omitempty"`
}

// ---------------------------------------------------------
// 2. Метаданные (Meta) - проверка перед открытием
// ---------------------------------------------------------

type SecretMetaResponse struct {
	HasPassword       bool `json:"hasPassword"`
	IsClientEncrypted bool `json:"isClientEncrypted"`
}

// ---------------------------------------------------------
// 3. Чтение записки (Read)
// ---------------------------------------------------------

type ReadSecretRequest struct {
	Password *string `json:"password,omitempty"`
}

type ReadSecretResponse struct {
	Text string `json:"text"`
}
