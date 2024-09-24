package main

import (
	"fmt"
)

// функция для вычисления длины строки
func stringlength(s string) int {
	return len(s)
}

func main() {
	var input string

	// ввод
	fmt.Print("Введите строку: ")
	fmt.Scanln(&input)

	// результат
	length := stringlength(input)
	fmt.Println("Длина строки:", length)
}
