package main

import (
	"fmt"
)

func main() {
	// создаем массив из 5 целых чисел
	array := [5]int{10, 20, 30, 40, 50}

	// создаем срез из массива
	slice := array[:]

	fmt.Println("Исходный срез:", slice)

	// добавление элемента в срез
	slice = append(slice, 60) // добавляем элемент 60
	fmt.Println("Срез после добавления 60:", slice)

	// удаление элемента из среза 
	indexToRemove := 4
	slice = append(slice[:indexToRemove], slice[indexToRemove+1:]...) // удаляем элемент с индексом 4 (то есть 50 т.к рассчет начинается от 0)
	fmt.Println("Срез после удаления элемента с индексом 5:", slice)

}
