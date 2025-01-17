package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"bytes"
)

func main() {
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1 - Зашифровать строку")
		fmt.Println("2 - Расшифровать строку")
		fmt.Println("3 - Выход")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			encryptString()
		case 2:
			decryptString()
		case 3:
			fmt.Println("Завершение программы.")
			return
		default:
			fmt.Println("Некорректный выбор, попробуйте снова.")
		}
	}
}

// Функция для шифрования данных
func encryptString() {
	fmt.Print("Введите строку для шифрования: ")
	var plaintext string
	fmt.Scanln(&plaintext)

	fmt.Print("Введите секретный ключ (16, 24 или 32 символа): ")
	var key string
	fmt.Scanln(&key)

	// Проверка длины ключа
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Некорректная длина ключа. Должен быть 16, 24 или 32 символа.")
		return
	}

	// Зашифровываем строку
	ciphertext, err := encrypt([]byte(plaintext), []byte(key))
	if err != nil {
		fmt.Println("Ошибка шифрования:", err)
		return
	}

	fmt.Println("Зашифрованная строка (Base64):", ciphertext)
}

// Функция для расшифровки данных
func decryptString() {
	fmt.Print("Введите строку для расшифровки (Base64): ")
	var ciphertext string
	fmt.Scanln(&ciphertext)

	fmt.Print("Введите секретный ключ (16, 24 или 32 символа): ")
	var key string
	fmt.Scanln(&key)

	// Проверка длины ключа
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Некорректная длина ключа. Должен быть 16, 24 или 32 символа.")
		return
	}

	// Дешифруем строку
	plaintext, err := decrypt(ciphertext, []byte(key))
	if err != nil {
		fmt.Println("Ошибка расшифровки:", err)
		return
	}

	fmt.Println("Расшифрованная строка:", plaintext)
}

// Шифрование данных с использованием AES (режим CBC)
func encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Добавляем паддинг, чтобы длина данных была кратна размеру блока AES (16 байт)
	plaintext = pad(plaintext, aes.BlockSize)

	// Генерация случайного IV (инициализационного вектора)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// Кодируем результат в Base64 для удобства
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Расшифровка данных с использованием AES (режим CBC)
func decrypt(ciphertext string, key []byte) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", errors.New("шифротекст слишком короткий")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	if len(ciphertextBytes)%aes.BlockSize != 0 {
		return "", errors.New("шифротекст некратен размеру блока")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	// Убираем паддинг
	plaintext, err := unpad(ciphertextBytes, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Функция для добавления паддинга (по стандарту PKCS#7)
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// Функция для удаления паддинга (по стандарту PKCS#7)
func unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, errors.New("длина данных некорректна для удаления паддинга")
	}

	padding := int(data[length-1])
	if padding > blockSize || padding == 0 {
		return nil, errors.New("некорректный паддинг")
	}

	return data[:length-padding], nil
}
