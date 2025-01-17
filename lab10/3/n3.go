package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	privateKeyFile = "private_key.pem"
	publicKeyFile  = "public_key.pem"
	signatureFile  = "signature.txt"
)

// Генерация ключей RSA и сохранение в файл
func generateKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Сохранение закрытого ключа
	privateFile, err := os.Create(privateKeyFile)
	if err != nil {
		return err
	}
	defer privateFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	if err := pem.Encode(privateFile, privatePem); err != nil {
		return err
	}

	// Сохранение открытого ключа
	publicKey := &privateKey.PublicKey
	publicFile, err := os.Create(publicKeyFile)
	if err != nil {
		return err
	}
	defer publicFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicPem := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return pem.Encode(publicFile, publicPem)
}

// Подпись сообщения с использованием закрытого ключа
func signMessage(message string) error {
	// Чтение закрытого ключа
	privatePem, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(privatePem)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	// Вычисление хэша сообщения
	hashed := sha256.Sum256([]byte(message))

	// Подпись
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return err
	}

	// Сохранение подписи в файл
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	return ioutil.WriteFile(signatureFile, []byte(signatureBase64), 0644)
}

// Проверка подписи
func verifySignature(message string) error {
	// Чтение открытого ключа
	publicPem, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(publicPem)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	rsaPublicKey := publicKey.(*rsa.PublicKey)

	// Чтение подписи из файла
	signatureBase64, err := ioutil.ReadFile(signatureFile)
	if err != nil {
		return err
	}
	signature, err := base64.StdEncoding.DecodeString(string(signatureBase64))
	if err != nil {
		return err
	}

	// Вычисление хэша сообщения
	hashed := sha256.Sum256([]byte(message))

	// Проверка подписи
	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return fmt.Errorf("подпись неверна: %v", err)
	}
	fmt.Println("Подпись подтверждена.")
	return nil
}

func main() {
	for {
		fmt.Println("Выберите действие:")
		fmt.Println("1 - Сгенерировать ключи")
		fmt.Println("2 - Подписать сообщение")
		fmt.Println("3 - Проверить подпись")
		fmt.Println("4 - Выход")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			err := generateKeys()
			if err != nil {
				fmt.Println("Ошибка генерации ключей:", err)
			} else {
				fmt.Println("Ключи успешно сгенерированы и сохранены.")
			}
		case 2:
			var message string
			fmt.Println("Введите сообщение для подписи:")
			fmt.Scan(&message)

			err := signMessage(message)
			if err != nil {
				fmt.Println("Ошибка при подписании:", err)
			} else {
				fmt.Println("Сообщение успешно подписано. Подпись сохранена в", signatureFile)
			}
		case 3:
			var message string
			fmt.Println("Введите сообщение для проверки подписи:")
			fmt.Scan(&message)

			err := verifySignature(message)
			if err != nil {
				fmt.Println("Ошибка при проверке подписи:", err)
			}
		case 4:
			fmt.Println("Завершение программы.")
			return
		default:
			fmt.Println("Некорректный выбор.")
		}
	}
}
