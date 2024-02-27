package utils

import "crypto/rand"

// GenerateKey - функция для генерации ключа шифрования
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256 требует  32 байта ключа
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
