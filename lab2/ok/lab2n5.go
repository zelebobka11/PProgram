package main

import (
	"fmt"
)

// структура Rectangle
type Rectangle struct {
	Width  float64
	Height float64
}

// метод для вычисления площади прямоугольника
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func main() {
	var width, height float64

	// ввод ширины и высоты прямоугольника 
	fmt.Print("Введите ширину прямоугольника: ")
	fmt.Scan(&width)
	fmt.Print("Введите высоту прямоугольника: ")
	fmt.Scan(&height)

	// экземпляр структуры Rectangle
	rect := Rectangle{Width: width, Height: height}

	// результат
	area := rect.Area()
	fmt.Printf("Площадь прямоугольника: %.2f\n", area)
}
