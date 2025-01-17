package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите идентификатор сессии (нажмите Enter, чтобы выполнить вход): ")
	sessionID, _ := reader.ReadString('\n')
	sessionID = strings.TrimSpace(sessionID)

	var csrfToken string

	if sessionID == "" {
		// Если идентификатор сессии не введён, выполняем вход
		fmt.Print("Введите имя пользователя (admin или user): ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Введите пароль: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		// Выполняем вход и получаем идентификатор сессии и CSRF-токен
		var err error
		sessionID, csrfToken, err = login(username, password)
		if err != nil {
			fmt.Println("Ошибка входа:", err)
			return
		}

		fmt.Println("Вход выполнен успешно.")
		fmt.Println("Ваш идентификатор сессии:", sessionID)
		fmt.Println("Ваш CSRF-токен:", csrfToken)
	} else {
		// Если идентификатор сессии введён, предполагаем, что CSRF-токен неизвестен
		fmt.Print("Введите CSRF-токен (нажмите Enter, если неизвестен): ")
		csrfToken, _ = reader.ReadString('\n')
		csrfToken = strings.TrimSpace(csrfToken)
	}

	// Тестируем доступ к маршруту /user
	fmt.Println("\nТестируем доступ к /user:")
	err := accessProtectedRoute("GET", "/user", sessionID, csrfToken)
	if err != nil {
		fmt.Println("Ошибка доступа к /user:", err)
	}

	// Тестируем доступ к маршруту /admin
	fmt.Println("\nТестируем доступ к /admin:")
	err = accessProtectedRoute("GET", "/admin", sessionID, csrfToken)
	if err != nil {
		fmt.Println("Ошибка доступа к /admin:", err)
	}

	// Тестируем POST-запрос к /update
	fmt.Println("\nТестируем POST-запрос к /update:")
	err = accessProtectedRoute("POST", "/update", sessionID, csrfToken)
	if err != nil {
		fmt.Println("Ошибка доступа к /update:", err)
	}

	// Тестируем POST-запрос к /update с неверным CSRF-токеном
	fmt.Println("\nТестируем POST-запрос к /update с неверным CSRF-токеном:")
	err = accessProtectedRoute("POST", "/update", sessionID, "invalid_csrf_token")
	if err != nil {
		fmt.Println("Ошибка доступа с неверным CSRF-токеном:", err)
	}
}

func login(username, password string) (string, string, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("http://localhost:8080/login?username=%s&password=%s", username, password)

	// Отправляем GET-запрос
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", "", fmt.Errorf("Ошибка аутентификации: %s", string(body))
	}

	// Читаем тело ответа (JSON)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Парсим JSON-ответ
	var response map[string]string
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", "", err
	}

	sessionID := response["session_id"]
	csrfToken := response["csrf_token"]

	return sessionID, csrfToken, nil
}

func accessProtectedRoute(method, route string, sessionID string, csrfToken string) error {
	// Формируем URL
	url := fmt.Sprintf("http://localhost:8080%s?session_id=%s", route, sessionID)

	// Создаём новый HTTP-запрос
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	// Если метод не GET или HEAD, добавляем CSRF-токен в заголовок
	if method != http.MethodGet && method != http.MethodHead {
		req.Header.Set("X-CSRF-Token", csrfToken)
	}

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Выводим статус-код и тело ответа
	fmt.Printf("Статус-код: %d\n", resp.StatusCode)
	fmt.Printf("Ответ: %s\n", string(body))

	return nil
}
