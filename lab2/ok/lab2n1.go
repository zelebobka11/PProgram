package main

import (
	"fmt"
)

func main() {
	var num int

	// ввод
	fmt.Print("Введите целое число: ")
	_, err := fmt.Scan(&num)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	// проверка на четность и нечетность
	if num%2 == 0 {
		fmt.Println("Число", num, "является четным.")
	} else {
		fmt.Println("Число", num, "является нечетным.")
	}
}
