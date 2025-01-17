package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
)

// Объект для обновления протокола до веб-сокета
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем подключение со всех источников
	},
}

var (
	clients    = make(map[*websocket.Conn]bool) // Подключенные клиенты
	broadcast  = make(chan Message)             // Канал для передачи сообщений
	clientsMux sync.Mutex                       // Мьютекс для защиты доступа к списку клиентов
)

// Структура для сообщения
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// Функция для обработки веб-сокет подключений
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновляем HTTP-соединение до веб-сокет соединения
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Добавляем клиента в список
	clientsMux.Lock()
	clients[ws] = true
	clientsMux.Unlock()

	// Постоянно читаем сообщения от клиента
	for {
		var msg Message
		// Чтение нового сообщения в формате JSON
		err := ws.ReadJSON(&msg)
		if err != nil {
			// Если клиент отключился, удаляем его из списка
			clientsMux.Lock()
			delete(clients, ws)
			clientsMux.Unlock()
			break
		}
		// Отправляем сообщение в канал
		broadcast <- msg
	}
}

// Функция для отправки сообщений всем клиентам
func handleMessages() {
	for {
		// Получаем сообщение из канала
		msg := <-broadcast

		// Рассылаем сообщение всем подключённым клиентам
		clientsMux.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				// Если произошла ошибка, закрываем соединение и удаляем клиента
				client.Close()
				delete(clients, client)
			}
		}
		clientsMux.Unlock()
	}
}

func main() {
	// Маршрут для подключения веб-сокетов
	http.HandleFunc("/ws", handleConnections)

	// Запуск горутины для рассылки сообщений
	go handleMessages()

	// Запуск HTTP сервера
	fmt.Println("Сервер запущен на http://localhost:5252")
	err := http.ListenAndServe(":5252", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
