package main

import (
	"fmt"
	"lab3/stringutils" // импорт пакета stringutils
)

func main() {
	var input string

	// ввод строки 
	fmt.Print("Введите строку для переворота: ")
	fmt.Scanln(&input)

	// вызов функции Reverse из пакета stringutils
	reversed := stringutils.Reverse(input)

	//результат
	fmt.Println("Перевернутая строка:", reversed)
}
