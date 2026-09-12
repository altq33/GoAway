package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"

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

// DecryptTextAES расшифровывает текст алгоритмом AES-256-GCM
func DecryptTextAES(encryptedHex, keyHex string) (string, error) {
	// 1. Декодируем hex-строки обратно в сырые байты
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", errors.New("неверный формат ключа")
	}
	encryptedBytes, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", errors.New("неверный формат зашифрованных данных")
	}

	// 2. Инициализируем шифр
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. Достаем Nonce
	nonceSize := aesGCM.NonceSize()
	if len(encryptedBytes) < nonceSize {
		return "", errors.New("поврежденные данные: слишком короткий текст")
	}

	// Отрезаем nonce (первые 12 байт) и сам зашифрованный текст (всё остальное)
	nonce := encryptedBytes[:nonceSize]
	ciphertext := encryptedBytes[nonceSize:]

	// 4. Расшифровываем (метод Open делает обратное методу Seal)
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("ошибка расшифровки: неверный ключ или данные повреждены")
	}

	return string(plaintext), nil
}
