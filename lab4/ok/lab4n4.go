package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// создаем новый сканер для считывания строки с ввода
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите строку: ")

	// считываем строку
	scanner.Scan()
	inputString := scanner.Text()

	// преобразуем строку в верхний регистр
	upperCaseString := strings.ToUpper(inputString)

	// выводим строку в верхнем регистре
	fmt.Println("Строка в верхнем регистре:", upperCaseString)
}
