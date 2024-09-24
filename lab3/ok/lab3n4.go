package main

import (
	"fmt"
)

func main() {
	// создаем массив из 5 целых чисел
	var numbers [5]int

	// заполняем массив
	for i := 0; i < len(numbers); i++ {
		fmt.Printf("Введите число для элемента %d: ", i)
		fmt.Scan(&numbers[i])
	}

	// вывод значений массива на экран
	fmt.Println("Введенные числа:")
	for _, number := range numbers {
		fmt.Println(number)
	}
}
