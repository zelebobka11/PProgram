package main

import (
	"fmt"
	"math/rand"
	"time"
)

// функция для вычисления факториала
func factorial(n int, resultChan chan int) {
	fmt.Printf("Горутина для факториала начала выполнение для числа %d\n", n)
	time.Sleep(2 * time.Second) // имитация задержки
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	resultChan <- result // отправляем результат в канал
	fmt.Printf("Факториал числа %d = %d\n", n, result)
}

// функция для генерации случайных чисел
func generateRandomNumbers(count int, resultChan chan []int) {
	fmt.Printf("Горутина для генерации случайных чисел начала выполнение для %d чисел\n", count)
	time.Sleep(1 * time.Second) // имитация задержки
	rand.Seed(time.Now().UnixNano())
	numbers := make([]int, count)
	for i := 0; i < count; i++ {
		numbers[i] = rand.Intn(100) // генерируем случайные числа от 0 до 100
	}
	resultChan <- numbers // отправляем результат в канал
	fmt.Printf("Случайные числа: %v\n", numbers)
}

// функция для вычисления суммы числового ряда
func sumSeries(n int, resultChan chan int) {
	fmt.Printf("Горутина для суммы ряда начала выполнение для числа %d\n", n)
	time.Sleep(3 * time.Second) // имитация задержки
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	resultChan <- sum // отправляем результат в канал
	fmt.Printf("Сумма числового ряда до %d = %d\n", n, sum)
}

func main() {
	// создаем каналы для получения результатов
	factorialChan := make(chan int)
	randomNumbersChan := make(chan []int)
	sumSeriesChan := make(chan int)

	// запускаем горутины
	go factorial(7, factorialChan)
	go generateRandomNumbers(5, randomNumbersChan)
	go sumSeries(52, sumSeriesChan)

	// получаем результаты из каналов
	factorialResult := <-factorialChan
	randomNumbersResult := <-randomNumbersChan
	sumSeriesResult := <-sumSeriesChan

	// выводим результаты
	fmt.Printf("Результат факториала: %d\n", factorialResult)
	fmt.Printf("Результат генерации случайных чисел: %v\n", randomNumbersResult)
	fmt.Printf("Результат суммы числового ряда: %d\n", sumSeriesResult)
}
