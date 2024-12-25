package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// Функция для сравнения введённого пользователем пароля с хешем реального пароля из бд
func ComparePbkdf2Sha256Hashes(storagePasswordHash string, enteredPassword string) bool {
	// Разбиваем полный хеш из бд на составные части
	params := strings.Split(storagePasswordHash, "$")
	if len(params) != 4 {
		log.Fatal("ComparePbkdf2Sha265Hashes: invalid hash format")
		return false
	}
	// Извлекаем соль хеша
	salt := params[2]
	// Извлекаем кол-во итераций хеширования
	iterations, err := strconv.Atoi(params[1])
	if err != nil {
		log.Fatalf("ComparePbkdf2Sha265Hashes: invalid hash format: %s", err.Error())
	}
	// Извлекаем закодированный хеш из бд и декодируем его из кодировки base64
	storeHash, errDecode := base64.StdEncoding.DecodeString(params[3])
	if errDecode != nil {
		log.Fatalf("ComparePbkdf2Sha265Hashes: unable to decode storage password hash %s", errDecode.Error())
	}

	// На основе извлечённых параметров хеша генерируем хеш для введённого пользователем пароля
	enteredPasswordHash := pbkdf2.Key([]byte(enteredPassword), []byte(salt), iterations, len(storeHash), sha256.New)

	// Срвниваем байтовые срезы хешей паролей
	return bytes.Equal(enteredPasswordHash, storeHash)
}

func GenerateRandomSalt() []byte {
	// Создаем байтовый срез под размер соли
	salt_size_bytes := 16 // 128 бит в django
	salt := make([]byte, salt_size_bytes)

	// Заполняем срез случайными числами
	_, err := rand.Read(salt[:])
	if err != nil {
		log.Printf("Unable to fill the salt %s", err.Error())
	}

	return salt
}

// Генерация хеша пароля алгоритмом pbkdf2_sha256,
// возврат результата в кодироовке base64 для хранения в бд
func GeneratePbkdf2Sha256Hash(password string) string {
	iterations := 870_000
	salt := GenerateRandomSalt()
	passwordBytes := []byte(password)
	algorithm := sha256.New
	keyLen := sha256.Size

	passwordHash := pbkdf2.Key(passwordBytes, salt, iterations, keyLen, algorithm)

	// Кодируем результат работы алгоритма pbkdf2 кодировкой base64 для хранения в бд
	passwordHashBase64 := base64.StdEncoding.EncodeToString(passwordHash)
	// кодируем соль для хранения вместе с хешем пароля (специально в формат совместимый с url)
	saltBase64 := base64.RawURLEncoding.EncodeToString(salt)

	passwordFormat := fmt.Sprintf("%s$%d$%s$%s", "pbkdf2_sha256", iterations, saltBase64, passwordHashBase64)

	return passwordFormat
}
