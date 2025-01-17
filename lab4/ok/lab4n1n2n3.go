package main

import (
	"fmt"
)

// функция для вычисления среднего возраста
func averageage(people map[string]int) float64 {
	if len(people) == 0 {
		return 0
	}

	totalAge := 0
	for _, age := range people {
		totalAge += age
	}
	return float64(totalAge) / float64(len(people))
}

// функция для удаления записи по имени
func deleteperson(people map[string]int, name string) {
	delete(people, name)
}

func main() {
	// создаем карту с именами людей и их возрастами
	people := map[string]int{
		"Shrek":    52,
		"Bebrich":  22,
		"Mujchina": 28,
	}

	// добавление нового человека
	people["Loshpedius"] = 44

	// вывод на экран
	fmt.Println("Все записи в карте:")
	for name, age := range people {
		fmt.Printf("Name: %s, Age: %d\n", name, age)
	}

	// средний возраст
	avgAge := averageage(people)
	fmt.Printf("Средний возраст: %.2f\n", avgAge)

	// удаление записи по заданному имени
	var nameToDelete string
	fmt.Print("Введите имя для удаления: ")
	fmt.Scan(&nameToDelete)
	deleteperson(people, nameToDelete)

	// выводим все записи на экран после удаления
	fmt.Println("Все записи в карте после удаления:")
	for name, age := range people {
		fmt.Printf("Name: %s, Age: %d\n", name, age)
	}

	// еще раз вычисляем и выводим средний возраст после удаления
	avgAge = averageage(people)
	fmt.Printf("Средний возраст после удаления: %.2f\n", avgAge)
}
