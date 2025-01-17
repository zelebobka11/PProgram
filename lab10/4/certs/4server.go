package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	select {
	case <-ctx.Done():
		fmt.Println("Сервер завершает соединение:", conn.RemoteAddr())
		return
	default:
	}

	// Чтение данных от клиента
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Сообщение от клиента:", message)

	// Ответ клиенту
	conn.Write([]byte("Сообщение получено!\n"))

	// Вывод на сервере подтверждения отправки
	fmt.Println("Сообщение отправлено клиенту:", conn.RemoteAddr())
}

func main() {
	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем канал для отслеживания прерываний (Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Загрузка сертификатов сервера
	cert, err := tls.LoadX509KeyPair("C:/Users/KLYKODIVER/Desktop/Pprogram/lab10/4/certs/server_cert.pem", "C:/Users/KLYKODIVER/Desktop/Pprogram/lab10/4/certs/server_key.pem")
	if err != nil {
		fmt.Println("Ошибка загрузки сертификата:", err)
		return
	}

	// Загрузка сертификата клиента (для взаимной аутентификации)
	clientCert, err := ioutil.ReadFile("C:/Users/KLYKODIVER/Desktop/Pprogram/lab10/4/certs/client_cert.pem")
	if err != nil {
		fmt.Println("Ошибка чтения сертификата клиента:", err)
		return
	}
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCert)

	// Настройка TLS конфигурации
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}

	// Запуск TLS-сервера
	ln, err := tls.Listen("tcp", ":5252", tlsConfig)
	if err != nil {
		fmt.Println("Ошибка запуска TLS сервера:", err)
		return
	}
	defer ln.Close()

	fmt.Println("TLS-сервер запущен на порту 5252")

	var wg sync.WaitGroup

	go func() {
		<-signalChan
		fmt.Println("\nПолучен сигнал для завершения работы сервера")
		cancel()        // Останавливаем новые соединения
		ln.Close()      // Закрываем слушатель
	}()

	// Основной цикл для принятия соединений
	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("Сервер завершает работу")
				wg.Wait() // Ждем завершения всех горутин
				return
			default:
				fmt.Println("Ошибка при принятии соединения:", err)
			}
			continue
		}

		wg.Add(1)
		go handleConnection(ctx, conn, &wg) // Обрабатываем соединение в горутине
	}
}
