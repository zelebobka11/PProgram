package main

import (
	"fmt"
)

// функцция для проверки числа
func checknumber(num int) string {
	if num > 0 {
		return "Positive"
	} else if num < 0 {
		return "Negative"
	} else {
		return "Zero"
	}
}

func main() {
	var num int

	// ввод числа
	fmt.Print("Введите число: ")
	_, err := fmt.Scan(&num)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	// результат
	result := checknumber(num)
	fmt.Println("Число является:", result)
}
