package main

import (
	"fmt"
	"time"
)

// Функция для вычисления суммы и разности двух чисел с плавающей запятой
func sumAndDiff(a, b float64) (float64, float64) {
	sum := a + b
	diff := a - b
	return sum, diff
}

// Функция для вычисления среднего значения трех чисел
func average(a, b, c float64) float64 {
	return (a + b + c) / 3
}

func main() {
	// 1. Вывод текущей даты и времени
	currentTime := time.Now()
	fmt.Println("Текущие дата и время:", currentTime)

	// 2. Создание переменных различных типов и вывод их на экран
	var i int = 42
	var f float64 = 3.14
	var s string = "Привет ворлд?!"
	var b bool = true
	fmt.Println("Переменные:")
	fmt.Println("int:", i)
	fmt.Println("float64:", f)
	fmt.Println("string:", s)
	fmt.Println("bool:", b)

	// 3. Использование краткой формы объявления переменных
	x := 10
	y := 20
	z := 2.5
	message := "Краткая форма объявления переменных"
	isGoFun := true
	fmt.Println("Краткая форма переменных:")
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)
	fmt.Println("message:", message)
	fmt.Println("isGoFun:", isGoFun)

	// 4. Арифметические операции с двумя целыми числами
	a := 15
	bb := 7
	fmt.Println("Арифметические операции:")
	fmt.Println("Сложение:", a, "+", bb, "=", a+bb)
	fmt.Println("Вычитание:", a, "-", bb, "=", a-bb)
	fmt.Println("Умножение:", a, "*", bb, "=", a*bb)
	fmt.Println("Деление:", a, "/", bb, "=", a/bb)
	fmt.Println("Остаток от деления:", a, "%", bb, "=", a%bb)

	// 5. Вызов функции для суммы и разности двух чисел с плавающей запятой
	num1 := 5.75
	num2 := 3.25
	sum, diff := sumAndDiff(num1, num2)
	fmt.Println("Сумма и разность чисел с плавающей запятой:")
	fmt.Println("Сумма:", sum)
	fmt.Println("Разность:", diff)

	// 6. Вычисление среднего значения трех чисел
	numA := 8.0
	numB := 6.0
	numC := 10.0
	avg := average(numA, numB, numC)
	fmt.Println("Среднее значение трех чисел:", avg)
}
