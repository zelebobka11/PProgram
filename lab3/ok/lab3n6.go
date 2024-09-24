package main

import (
	"fmt"
)

func main() {
	// создаем срез строк
	strings := []string{
		"туз",
		"король",
		"валет",
		"дама",
		"шестерочка",
	}

	// инициализируем переменные для хранения самой длинной строки и её длины
	longeststring := ""
	maxlength := 0

	// ищем самую длинную строку
	for _, str := range strings {
		if len(str) > maxlength {
			longeststring = str
			maxlength = len(str)
		}
	}

	// результат
	fmt.Println("Самая длинная строка:", longeststring)
	fmt.Println("Длина самой длинной строки:", maxlength)
}
