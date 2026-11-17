package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomEmailLength     = 12
	randomPasswordLength  = 12
	passwordRequiredChars = "Aa1"
)

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		randCharsetIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(fmt.Sprintf("Не удалось сгенерировать случайное число: %s", err.Error()))
		}
		b[i] = charset[randCharsetIndex.Int64()]
	}

	return string(b)
}

// RandomEmail возвращает уникальный email для регистрации тестового пользователя
func RandomEmail() string {
	return fmt.Sprintf("autotest_%s@example.com", strings.ToLower(RandomString(randomEmailLength)))
}

// RandomPassword возвращает пароль по правилам сервиса: заглавные и строчные буквы, цифры, не короче 8 символов
func RandomPassword() string {
	return passwordRequiredChars + RandomString(randomPasswordLength)
}
