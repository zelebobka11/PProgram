package main

import (
	"fmt"
)

func main() {
	var n, sum int 

	// запрос на количество чисел
	fmt.Print("Введите количество чисел: ")
	fmt.Scan(&n) 

	// срез для хранения чисел
	numbers := make([]int, n)

	fmt.Println("Введите числа:")
	// цикл для ввода чисел
	for i := 0; i < n; i++ {
		fmt.Scan(&numbers[i]) 
		sum += numbers[i]     
	}

	// результат
	fmt.Printf("Сумма введенных чисел: %d\n", sum)
}
