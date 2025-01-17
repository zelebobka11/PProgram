package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Структура для представления сессии пользователя
type Session struct {
	Username  string
	Role      string
	Expires   time.Time
	CSRFToken string
}

// Хранилище сессий в памяти
var sessions = make(map[string]Session)
var sessionsFile = "sessions.json"
var mu sync.Mutex // Для синхронизации доступа к сессиям

// Генерация уникального идентификатора сессии
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Генерация CSRF-токена
func generateCSRFToken() (string, error) {
	return generateSessionID()
}

// Загрузка сессий из файла
func loadSessions() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(sessionsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Файл не существует, ничего не делаем
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sessions)
	if err != nil {
		return err
	}
	return nil
}

// Сохранение сессий в файл
func saveSessions() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create(sessionsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(sessions)
	if err != nil {
		return err
	}
	return nil
}

// Аутентификация пользователя и создание сессии
func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	// Пример статической аутентификации
	if (username == "admin" || username == "user") && password == "password" {
		role := username

		// Генерация идентификатора сессии
		sessionID, err := generateSessionID()
		if err != nil {
			http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
			return
		}

		// Генерация CSRF-токена
		csrfToken, err := generateCSRFToken()
		if err != nil {
			http.Error(w, "Ошибка создания CSRF-токена", http.StatusInternalServerError)
			return
		}

		// Установка времени истечения сессии
		expirationTime := time.Now().Add(15 * time.Minute)

		// Создание и сохранение сессии
		session := Session{
			Username:  username,
			Role:      role,
			Expires:   expirationTime,
			CSRFToken: csrfToken,
		}

		mu.Lock()
		sessions[sessionID] = session
		mu.Unlock()

		// Сохранение сессий в файл
		err = saveSessions()
		if err != nil {
			http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
			return
		}

		// Отправка идентификатора сессии и CSRF-токена клиенту
		response := map[string]string{
			"session_id": sessionID,
			"csrf_token": csrfToken,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
	}
}

// Middleware для проверки сессии и ролей
func authorize(next http.HandlerFunc, requiredRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение идентификатора сессии из заголовка или параметра
		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			sessionID = r.Header.Get("Session-Id")
			if sessionID == "" {
				http.Error(w, "Сессия не найдена", http.StatusUnauthorized)
				return
			}
		}

		// Загрузка сессий из файла
		err := loadSessions()
		if err != nil {
			http.Error(w, "Ошибка загрузки сессий", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		session, exists := sessions[sessionID]
		mu.Unlock()

		if !exists {
			http.Error(w, "Сессия не найдена", http.StatusUnauthorized)
			return
		}

		// Проверка истечения сессии
		if session.Expires.Before(time.Now()) {
			mu.Lock()
			delete(sessions, sessionID)
			mu.Unlock()
			saveSessions()
			http.Error(w, "Сессия истекла", http.StatusUnauthorized)
			return
		}

		// Проверка роли пользователя
		if session.Role != requiredRole {
			http.Error(w, "Недостаточно прав", http.StatusForbidden)
			return
		}

		// Защита от CSRF для небезопасных методов
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			csrfToken := r.Header.Get("X-CSRF-Token")
			if csrfToken == "" || csrfToken != session.CSRFToken {
				http.Error(w, "Неверный CSRF-токен", http.StatusForbidden)
				return
			}
		}

		// Вызов следующего обработчика
		next.ServeHTTP(w, r)
	}
}

// Обработчик для администраторов
func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Доступ к данным только для администраторов")
}

// Обработчик для пользователей
func userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Доступ к данным для пользователей")
}

// Обработчик для обновления данных (POST-запрос)
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Обновление данных пользователя успешно")
}

func main() {
	// Загрузка сессий из файла при запуске сервера
	err := loadSessions()
	if err != nil {
		log.Fatal("Ошибка загрузки сессий:", err)
	}

	http.HandleFunc("/login", loginHandler)

	// Защищённые маршруты
	http.HandleFunc("/admin", authorize(adminHandler, "admin"))
	http.HandleFunc("/user", authorize(userHandler, "user"))
	http.HandleFunc("/update", authorize(updateUserHandler, "user"))

	fmt.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
