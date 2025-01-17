package main

import (
	"fmt"
)

// определение структуры Person с полями name и age
type Person struct {
	name string
	age  int
}

// метод для вывода информации о человеке
func (p Person) printInfo() {
	fmt.Printf("Имя: %s, Возраст: %d лет\n", p.name, p.age)
}

// метод для увеличения возраста на год
func (p *Person) birthday() {
	p.age += 1
	fmt.Printf("Харош, %s! Теперь ты старый!Ведь тебе: %d \n", p.name, p.age)
}

func main() {
	// создание экземпляра структуры Person
	person := Person{name: "Бебра", age: 51}

	// выводим информацию о человеке
	person.printInfo()

	// увеличиваем возраст на год
	person.birthday()

	// выводим обновленную информацию
	person.printInfo()
}
