package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Подключаемся к TCP-серверу на localhost:5252
	conn, err := net.Dial("tcp", "localhost:5252")
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
