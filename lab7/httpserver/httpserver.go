package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Структура для представления данных, отправленных в POST запросе
type Data struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Middleware для логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Логируем метод и URL запроса
		log.Printf("Запрос: %s %s", r.Method, r.URL.Path)
		// Выполняем следующий обработчик
		next.ServeHTTP(w, r)
		// Логируем время выполнения
		log.Printf("Время выполнения: %v\n", time.Since(start))
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Обработка GET запроса /hello
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Привет, добро пожаловать на Go HTTP сервер!"))
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// Обработка POST запроса /data
	if r.Method == http.MethodPost {
		// Чтение тела запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
			return
		}

		// Парсинг JSON данных
		var data Data
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		// Вывод содержимого JSON в консоль
		fmt.Printf("Полученные данные: Имя=%s, Email=%s\n", data.Name, data.Email)

		// Ответ клиенту
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "получено"}`))
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Маршрутизация
	mux := http.NewServeMux()

	// Регистрируем маршруты
	mux.HandleFunc("/hello", helloHandler) // GET /hello
	mux.HandleFunc("/data", dataHandler)   // POST /data

	// Оборачиваем маршрутизатор в middleware
	loggedMux := loggingMiddleware(mux)

	// Запуск сервера на порту 5252
	fmt.Println("Запуск сервера на http://localhost:5252...")
	err := http.ListenAndServe(":5252", loggedMux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

//              curl http://localhost:5252/hello
//              Invoke-WebRequest -Uri http://localhost:5252/data -Method POST -Body '{"name": "Michael", "email": "miha@example.com"}' -Headers @{"Content-Type"="application/json"}
