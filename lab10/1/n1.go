package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1 - Вычислить хэш строки")
		fmt.Println("2 - Проверить целостность данных")
		fmt.Println("3 - Выход")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Ошибка ввода. Попробуйте снова.")
			continue
		}

		switch choice {
		case 1:
			calculateHash()
		case 2:
			checkDataIntegrity()
		case 3:
			fmt.Println("Завершение программы.")
			return
		default:
			fmt.Println("Некорректный выбор, попробуйте снова.")
		}
	}
}

// Выбор алгоритма хэширования
func selectHashAlgorithm() (hash.Hash, string, error) {
	fmt.Println("Выберите хэш-функцию:")
	fmt.Println("1 - MD5")
	fmt.Println("2 - SHA-256")
	fmt.Println("3 - SHA-512")

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return nil, "", errors.New("Некорректный ввод хэш-функции")
	}

	switch choice {
	case 1:
		return md5.New(), "MD5", nil
	case 2:
		return sha256.New(), "SHA-256", nil
	case 3:
		return sha512.New(), "SHA-512", nil
	default:
		return nil, "", errors.New("Некорректный выбор хэш-функции")
	}
}

// Вычисление хэша строки
func calculateHash() {
	// Ввод строки
	fmt.Print("Введите строку для хэширования: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Удаление символов новой строки

	// Выбор алгоритма хэширования
	hasher, algoName, err := selectHashAlgorithm()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Вычисление хэша
	hasher.Write([]byte(input))
	hashValue := hex.EncodeToString(hasher.Sum(nil))

	fmt.Printf("Алгоритм: %s\n", algoName)
	fmt.Printf("Хэш: %s\n", hashValue)
}

// Проверка целостности данных
func checkDataIntegrity() {
	// Ввод строки для проверки
	fmt.Print("Введите строку для проверки: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Удаление символов новой строки

	// Ввод хэша для проверки
	fmt.Print("Введите хэш для проверки: ")
	inputHash, _ := reader.ReadString('\n')
	inputHash = strings.TrimSpace(inputHash) // Удаление символов новой строки

	// Выбор алгоритма хэширования
	hasher, algoName, err := selectHashAlgorithm()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Вычисление хэша строки
	hasher.Write([]byte(input))
	computedHash := hex.EncodeToString(hasher.Sum(nil))

	// Сравнение введённого хэша с вычисленным
	if strings.EqualFold(computedHash, inputHash) {
		fmt.Printf("Хэши совпадают. Строка и её хэш соответствуют. (Алгоритм: %s)\n", algoName)
	} else {
		fmt.Printf("Хэши НЕ совпадают. Строка и её хэш не соответствуют. (Алгоритм: %s)\n", algoName)
	}
}
