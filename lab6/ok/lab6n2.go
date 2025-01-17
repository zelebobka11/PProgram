package main

import (
	"fmt"
)

// функция для генерации чисел фибоначчи и отправки их в канал
func generateFibonacci(n int, ch chan int) {
	// инициализация первых двух чисел фибоначчи
	a, b := 0, 1
	for i := 0; i < n; i++ {
		ch <- a // отправляем значение в канал
		a, b = b, a+b
	}
	close(ch) // закрывает канал и сигнализирует о прекращении поступления данных
}

// функция для чтения из канала и вывода значений
func printFibonacci(ch chan int) {
	// чтение из канала до его закрытия
	for num := range ch {
		fmt.Println(num) // вывод значений на экран
	}
}

func main() {
	// создание канала для передачи чисел фибоначчи
	fibChan := make(chan int)

	// запуск горутины для генерации чисел фибоначчи
	go generateFibonacci(10, fibChan)

	// запуск горутины для чтения и вывода чисел
	printFibonacci(fibChan)
}
