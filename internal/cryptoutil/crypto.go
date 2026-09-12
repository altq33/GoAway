package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// GenerateHexID генерирует криптографически стойкую случайную строку (hex)
func GenerateHexID(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword создает bcrypt-хэш пароля и возвращает указатель на строку
func HashPassword(password string) (*string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashStr := string(hashBytes)
	return &hashStr, nil
}

// EncryptTextAES шифрует текст алгоритмом AES-256-GCM.
// Возвращает (зашифрованный текст в hex, сгенерированный ключ в hex, ошибку)
func EncryptTextAES(text string) (string, string, error) {
	// 1. Генерируем 32 байта для ключа AES-256
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", "", err
	}

	// 2. Инициализируем шифр
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	// 3. Генерируем Nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", "", err
	}

	// 4. Шифруем и приклеиваем Nonce в начало
	encryptedBytes := aesGCM.Seal(nonce, nonce, []byte(text), nil)

	// Возвращаем (текст, ключ, nil)
	return hex.EncodeToString(encryptedBytes), hex.EncodeToString(key), nil
}
