package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// Загрузка сертификата сервера для верификации
	serverCert, err := ioutil.ReadFile("C:/Users/KLYKODIVER/Desktop/Pprogram/lab10/4/certs/ca_cert.pem")
	if err != nil {
		fmt.Println("Ошибка чтения сертификата сервера:", err)
		return
	}
	serverCertPool := x509.NewCertPool()
	serverCertPool.AppendCertsFromPEM(serverCert)

	// Загрузка сертификата клиента для аутентификации
	clientCert, err := tls.LoadX509KeyPair("C:/Users/KLYKODIVER/Desktop/Pprogram/lab10/4/certs/client_cert.pem", "C:/Users/KLYKODIVER/Desktop/Pprogram/lab10/4/certs/client_key.pem")
	if err != nil {
		fmt.Println("Ошибка загрузки сертификата клиента:", err)
		return
	}

	// Настройка TLS конфигурации для клиента
	tlsConfig := &tls.Config{
		RootCAs:      serverCertPool,  // Добавляем корневой сертификат сервера
		Certificates: []tls.Certificate{clientCert}, // Добавляем сертификат клиента
	}

	// Подключаемся к TLS-серверу на localhost:5252
	conn, err := tls.Dial("tcp", "localhost:5252", tlsConfig)
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Ввод сообщения от пользователя
	fmt.Print("Введите сообщение для отправки серверу: ")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')

	// Отправляем сообщение серверу
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}

	// Чтение ответа от сервера
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при получении ответа:", err)
		return
	}

	// Выводим ответ от сервера
	fmt.Printf("Ответ от сервера: %s", response)
}
