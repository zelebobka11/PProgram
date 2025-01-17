package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	Username string
	Expiry   time.Time
}

var (
	users    = make(map[string]string) // Хранение пользователей: username -> password
	sessions = make(map[string]Session)
	mu       sync.Mutex // Мьютекс для синхронизации доступа
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Проверяем, существует ли пользователь
	if _, exists := users[creds.Username]; exists {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	// Сохраняем пользователя
	users[creds.Username] = creds.Password
	fmt.Printf("Зарегистрирован новый пользователь: %s\n", creds.Username)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Регистрация успешна"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Проверяем, существует ли пользователь и совпадает ли пароль
	if password, exists := users[creds.Username]; !exists || password != creds.Password {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	// Генерируем токен сессии
	token := generateSessionToken()
	expiry := time.Now().Add(30 * time.Minute)
	sessions[token] = Session{
		Username: creds.Username,
		Expiry:   expiry,
	}

	// Отправляем токен клиенту
	resp := map[string]string{"token": token}
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func generateSessionToken() string {
	const tokenLength = 32
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, tokenLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Middleware для проверки авторизации
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < len("Bearer ") {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}
		token := authHeader[len("Bearer "):]

		mu.Lock()
		session, exists := sessions[token]
		mu.Unlock()

		if !exists || session.Expiry.Before(time.Now()) {
			http.Error(w, "Сессия истекла или недействительна", http.StatusUnauthorized)
			return
		}

		// Обновляем время окончания сессии
		mu.Lock()
		session.Expiry = time.Now().Add(30 * time.Minute)
		sessions[token] = session
		mu.Unlock()

		// Добавляем имя пользователя в контекст запроса (если необходимо)
		next(w, r)
	}
}

// Обработчики для CRUD операций с пользователями

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var (
	usersData = make(map[int]User) // Хранение данных пользователей: ID -> User
	userIDSeq = 1                  // Счетчик для генерации ID пользователей
)

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user.ID = userIDSeq
	usersData[userIDSeq] = user
	userIDSeq++
	mu.Unlock()

	fmt.Printf("Добавлен пользователь: %s\n", user.Username)
	w.WriteHeader(http.StatusCreated)
	jsonResp, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/users/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	mu.Lock()
	defer mu.Unlock()

	if _, exists := usersData[id]; !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	delete(usersData, id)
	fmt.Printf("Удален пользователь с ID: %d\n", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь удален успешно"))
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/users/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := usersData[id]; !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	user.ID = id
	usersData[id] = user
	fmt.Printf("Обновлен пользователь с ID: %d\n", id)
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func viewUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/users/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	mu.Lock()
	user, exists := usersData[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	jsonResp, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func viewAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	usersList := make([]User, 0, len(usersData))
	for _, user := range usersData {
		usersList = append(usersList, user)
	}

	jsonResp, _ := json.Marshal(usersList)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authMiddleware(viewAllUsersHandler)(w, r)
		case http.MethodPost:
			authMiddleware(addUserHandler)(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authMiddleware(viewUserHandler)(w, r)
		case http.MethodPut:
			authMiddleware(updateUserHandler)(w, r)
		case http.MethodDelete:
			authMiddleware(deleteUserHandler)(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Сервер запущен на порту 5252...")
	log.Fatal(http.ListenAndServe(":5252", nil))
}
