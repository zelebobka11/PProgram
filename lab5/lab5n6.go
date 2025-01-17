package main

import (
	"fmt"
)

type Stringer interface {
	String() string
}

// сСтруктура Book содержит информацию о книге: название, автор и количество страниц
type Book struct {
	Title  string
	Author string
	Pages  int
}

// метод String для структуры Book
func (b Book) String() string {
	return fmt.Sprintf("Title: %s, Author: %s, Pages: %d", b.Title, b.Author, b.Pages)
}

func main() {
	// создание экземпляра структуры Book с данными о книге
	book := Book{
		Title:  "Что в мешочке", // название
		Author: "Сигма",         // автор
		Pages:  52,              // колво страниц
	}

	// результат
	fmt.Println(book)

}
