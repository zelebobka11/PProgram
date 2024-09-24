package main

import (
	"fmt"
	"lab3/mathutils" //иммпорт пакета mathutils
)

func main() {
	var num int

	// ввод числа 
	fmt.Print("Введите число для вычисления факториала: ")
	fmt.Scan(&num)

	// вызов функции из пакета mathutils
	result := mathutils.Factorial(num)

	if result == -1 {
		fmt.Println("Ошибка: факториал для отрицательного числа не определен.")
	} else {
		fmt.Printf("Факториал числа %d равен: %d\n", num, result)
	}
}
