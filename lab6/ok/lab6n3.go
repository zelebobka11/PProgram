package main

import (
	"fmt"
	"math/rand"
	"time"
)

// горутина для генерации случайных чисел
func generateRandomNumbers(numbersChan chan int) {
	for {
		num := rand.Intn(100) // генерация случайного числа от 0 до 99
		numbersChan <- num    // отправка числа в канал
		time.Sleep(500 * time.Millisecond) // имитация задержки (по приколу)
	}
}

// горутина для определения четности/нечётности чисел
func checkEvenOdd(numbersChan chan int, messagesChan chan string) {
	for num := range numbersChan {
		if num%2 == 0 {
			messagesChan <- fmt.Sprintf("%d Четное", num) // четное число
		} else {
			messagesChan <- fmt.Sprintf("%d Нечетное", num)  // нечетное число
		}
	}
}

func main() {
	// инициализация каналов
	numbersChan := make(chan int)
	messagesChan := make(chan string)

	// запуск горутин
	go generateRandomNumbers(numbersChan)
	go checkEvenOdd(numbersChan, messagesChan)

	// счётчик для отслеживания числа обработанных результатов
	resultCount := 0
	const maxResults = 6 // максимальное количество результатов

	// цикл с использованием select для управления каналами
	for {
		select {
		case message := <-messagesChan:
			fmt.Println(message) // вывод сообщения о четности/нечетности числа
			resultCount++       
		case num := <-numbersChan:
			fmt.Printf("Received number: %d\n", num) // вывод самого числа
			resultCount++                            
		}

		if resultCount >= maxResults {
			fmt.Println("Конец...")
			return 
		}
	}
}
