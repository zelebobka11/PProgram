package main

import (
	"fmt"
)

// функция для вычисления среднего значения двух целых чисел
func average(a int, b int) float64 {
	return float64(a+b) / 2.0
}

func main() {
	var num1, num2 int

	// ввод двух целых чисел
	fmt.Print("Введите первое целое число: ")
	fmt.Scan(&num1)
	fmt.Print("Введите второе целое число: ")
	fmt.Scan(&num2)

	// результат
	avg := average(num1, num2)
	fmt.Printf("Среднее значение чисел %d и %d равно: %.2f\n", num1, num2, avg)
}
