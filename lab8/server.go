package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite" // Используем драйвер SQLite
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

// Централизованная обработка ошибок
func handleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

// Инициализация базы данных
func initDB() {
	var err error
	// Открытие базы данных SQLite
	db, err = sql.Open("sqlite", "./users.db") // Используем "sqlite" для modernc.org/sqlite
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы пользователей, если её нет
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

// Валидация данных пользователя
func validateUser(user User) error {
	if user.Name == "" {
		return errors.New("имя не может быть пустым")
	}
	if user.Age <= 0 {
		return errors.New("возраст должен быть положительным числом")
	}
	return nil
}

// Получение списка пользователей с пагинацией и фильтрацией
func getUsers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	nameFilter := r.URL.Query().Get("name")
	ageFilter := r.URL.Query().Get("age")

	query := "SELECT id, name, age FROM users WHERE 1=1"
	var params []interface{}

	// Добавляем фильтрацию по имени
	if nameFilter != "" {
		query += " AND name LIKE ?"
		params = append(params, "%"+nameFilter+"%")
	}

	// Добавляем фильтрацию по возрасту
	if ageFilter != "" {
		query += " AND age = ?"
		age, err := strconv.Atoi(ageFilter)
		if err == nil {
			params = append(params, age)
		}
	}

	query += " LIMIT ? OFFSET ?"
	params = append(params, limit, offset)

	rows, err := db.Query(query, params...)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Получение пользователя по ID
func getUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(w, errors.New("некорректный ID пользователя"), http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			handleError(w, errors.New("пользователь не найден"), http.StatusNotFound)
		} else {
			handleError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Создание нового пользователя
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		handleError(w, errors.New("некорректный формат данных"), http.StatusBadRequest)
		return
	}

	// Валидация данных
	err = validateUser(user)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	user.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Обновление информации о пользователе
func updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(w, errors.New("некорректный ID пользователя"), http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		handleError(w, errors.New("некорректный формат данных"), http.StatusBadRequest)
		return
	}

	// Валидация данных
	err = validateUser(user)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", user.Name, user.Age, id)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	user.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Удаление пользователя
func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(w, errors.New("некорректный ID пользователя"), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsers(w, r)
		} else if r.Method == http.MethodPost {
			createUser(w, r)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUserByID(w, r)
		} else if r.Method == http.MethodPut {
			updateUser(w, r)
		} else if r.Method == http.MethodDelete {
			deleteUser(w, r)
		}
	})

	fmt.Println("Сервер запущен на порту 5252...")
	log.Fatal(http.ListenAndServe(":5252", nil))
}



/*
Тестирование программы:
1. Получение списка пользователей (GET /users):
   curl http://localhost:5252/users

2. Получение информации о конкретном пользователе (GET /users/{id}):
	curl http://localhost:5252/users/1
	curl http://localhost:5252/users?age=20
	curl http://localhost:5252/users?name=Mkas


3. Добавление нового пользователя (POST /users):
   curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"Mkas\", \"age\":20}" http://localhost:5252/users

4. Обновление информации о пользователе (PUT /users/{id}):
   curl -X PUT -H "Content-Type: application/json" -d "{\"name\":\"Mkas Updated\", \"age\":40}" http://localhost:5252/users/1

5. Удаление пользователя (DELETE /users/{id}):
    curl -X DELETE http://localhost:5252/users/1

*/