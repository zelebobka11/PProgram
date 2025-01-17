package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const serverURL = "http://localhost:5252"

var sessionToken string

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n--- Меню ---")
		fmt.Println("1. Регистрация")
		fmt.Println("2. Вход")
		fmt.Println("3. Добавить пользователя")
		fmt.Println("4. Удалить пользователя")
		fmt.Println("5. Обновить информацию пользователя")
		fmt.Println("6. Просмотреть пользователя")
		fmt.Println("7. Выйти из аккаунта")
		fmt.Println("8. Выход")
		fmt.Println("9. Просмотреть всех пользователей")
		fmt.Print("Выберите опцию: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			register()
		case "2":
			login()
		case "3":
			addUser()
		case "4":
			deleteUser()
		case "5":
			updateUser()
		case "6":
			viewUser()
		case "7":
			logout()
		case "8":
			fmt.Println("Выход из программы.")
			return
		case "9":
			viewAllUsers()
		default:
			fmt.Println("Неверный выбор, попробуйте снова.")
		}
	}
}

func register() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите имя пользователя: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Введите пароль: ")
	password, _ := reader.ReadString('\n')

	creds := Credentials{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	data, _ := json.Marshal(creds)
	resp, err := http.Post(serverURL+"/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при регистрации:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка регистрации:", string(body))
		return
	}

	fmt.Println("Регистрация успешна.")
}

func login() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите имя пользователя: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Введите пароль: ")
	password, _ := reader.ReadString('\n')

	creds := Credentials{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	data, _ := json.Marshal(creds)
	resp, err := http.Post(serverURL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при входе:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка входа:", string(body))
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)
	sessionToken = result["token"]

	fmt.Println("Вход выполнен успешно.")
}

func addUser() {
	if sessionToken == "" {
		fmt.Println("Пожалуйста, войдите в систему.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите имя пользователя: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Введите email: ")
	email, _ := reader.ReadString('\n')

	user := User{
		Username: strings.TrimSpace(username),
		Email:    strings.TrimSpace(email),
	}

	data, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", serverURL+"/users", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при добавлении пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка добавления пользователя:", string(body))
		return
	}

	// Обработка ответа сервера
	body, _ := ioutil.ReadAll(resp.Body)
	var addedUser User
	json.Unmarshal(body, &addedUser)

	fmt.Println("Пользователь добавлен успешно.")
	fmt.Printf("ID нового пользователя: %d\n", addedUser.ID)
}

func deleteUser() {
	if sessionToken == "" {
		fmt.Println("Пожалуйста, войдите в систему.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите ID пользователя для удаления: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	req, _ := http.NewRequest("DELETE", serverURL+"/users/"+idStr, nil)
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при удалении пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка удаления пользователя:", string(body))
		return
	}

	fmt.Println("Пользователь удален успешно.")
}

func updateUser() {
	if sessionToken == "" {
		fmt.Println("Пожалуйста, войдите в систему.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите ID пользователя для обновления: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	fmt.Print("Введите новое имя пользователя: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Введите новый email: ")
	email, _ := reader.ReadString('\n')

	user := User{
		Username: strings.TrimSpace(username),
		Email:    strings.TrimSpace(email),
	}

	data, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", serverURL+"/users/"+idStr, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при обновлении пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка обновления пользователя:", string(body))
		return
	}

	fmt.Println("Пользователь обновлен успешно.")
}

func viewUser() {
	if sessionToken == "" {
		fmt.Println("Пожалуйста, войдите в систему.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите ID пользователя для просмотра: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	req, _ := http.NewRequest("GET", serverURL+"/users/"+idStr, nil)
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при получении пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка получения пользователя:", string(body))
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var user User
	json.Unmarshal(body, &user)

	fmt.Println("Информация о пользователе:")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("Имя пользователя: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
}

func viewAllUsers() {
	if sessionToken == "" {
		fmt.Println("Пожалуйста, войдите в систему.")
		return
	}

	req, _ := http.NewRequest("GET", serverURL+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при получении списка пользователей:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка получения списка пользователей:", string(body))
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var users []User
	json.Unmarshal(body, &users)

	fmt.Println("Список пользователей:")
	for _, user := range users {
		fmt.Printf("ID: %d, Имя пользователя: %s, Email: %s\n", user.ID, user.Username, user.Email)
	}
}

func logout() {
	sessionToken = ""
	fmt.Println("Вы вышли из аккаунта.")
}
