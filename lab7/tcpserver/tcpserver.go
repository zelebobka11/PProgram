package main

import (
	"bufio"
	"context"
	"fmt"
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

	// Запуск сервера на порту 5252
	ln, err := net.Listen("tcp", ":5252")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Сервер запущен на порту 5252")

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
