package main

import (
	"fmt"
)

func main() {
	var n int

	// запрсо на количество элементов массива
	fmt.Print("Введите количество чисел: ")
	fmt.Scan(&n)

	// срез для хранения чисел
	numbers := make([]int, n)

	// ввод чисел
	fmt.Println("Введите числа:")
	for i := 0; i < n; i++ {
		fmt.Scan(&numbers[i])
	}

	// вывод чисел в обратном порядке
	fmt.Println("Числа в обратном порядке:")
	for i := n - 1; i >= 0; i-- { // цикл от конца к началу
		fmt.Printf("%d ", numbers[i])
	}
}
