package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// Функция для спавнения введённого пользователем пароля с хешем реального пароля из бд\
func ComparePbkdf2Sha265Hashes(storagePasswordHash string, enteredPassword string) bool {
	// Разбиваем полный хеш из бд на составные части
	params := strings.Split(storagePasswordHash, "$")
	if len(params) != 4 {
		log.Fatal("ComparePbkdf2Sha265Hashes: invalid hash format")
		return false
	}
	// Извлека5ем соль хеша
	salt := params[2]
	// Извлекаем кол-во итераций хеширования
	iterations, err := strconv.Atoi(params[1])
	if err != nil {
		log.Fatalf("ComparePbkdf2Sha265Hashes: invalid hash format: %s", err.Error())
	}
	// Извлекаем закодированный хеш из бд и декодируем его из кодировки base64
	storeHash, errDecode := base64.StdEncoding.DecodeString(params[3])
	if errDecode != nil {
		log.Fatalf("ComparePbkdf2Sha265Hashes: unable to decode starage password hash %s", errDecode.Error())
	}

	// На основе извлечённых параметров хеша генерируем хеш для введённого пользователем пароля
	enteredPasswordHash := pbkdf2.Key([]byte(enteredPassword), []byte(salt), iterations, len(storeHash), sha256.New)

	// Срвниваем байтовые срезы хешей паролей
	return bytes.Equal(enteredPasswordHash, storeHash)
}
