package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Настройки ключей
var jwtKey = []byte("my_secret_key")
var store = sessions.NewCookieStore([]byte("session_secret_key"))

// Структура для хранения данных пользователя
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// Генерация CSRF токена
func generateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Маршрут для входа и получения JWT токена
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Пример статичной аутентификации
	var role string
	if creds.Username == "admin" && creds.Password == "password" {
		role = "admin"
	} else if creds.Username == "user" && creds.Password == "password" {
		role = "user"
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Генерация CSRF токена
	csrfToken, err := generateCSRFToken()
	if err != nil {
		http.Error(w, "CSRF Token Error", http.StatusInternalServerError)
		return
	}

	// Создание сессии
	session, _ := store.Get(r, "session")
	session.Values["username"] = creds.Username
	session.Values["role"] = role
	session.Values["csrfToken"] = csrfToken
	session.Save(r, w)

	// Отправка токенов клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "csrfToken": csrfToken})
}

// Middleware для проверки JWT и CSRF токена
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка JWT токена
		tokenString := r.Header.Get("Authorization")
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Проверка роли
		role := claims.Role
		if role != "admin" && r.URL.Path == "/admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Проверка CSRF токена
		session, _ := store.Get(r, "session")
		csrfToken, ok := session.Values["csrfToken"].(string)
		if !ok || r.Header.Get("X-CSRF-Token") != csrfToken {
			http.Error(w, "Invalid CSRF Token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Защищенные маршруты
func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Admin!")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, User!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")

	// Защищенные маршруты
	r.Handle("/admin", authMiddleware(http.HandlerFunc(adminHandler))).Methods("GET")
	r.Handle("/user", authMiddleware(http.HandlerFunc(userHandler))).Methods("GET")

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", r)
}
