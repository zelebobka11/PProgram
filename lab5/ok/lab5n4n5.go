package main

import (
	"fmt"
	"math"
)

// интерфейс Shape
type Shape interface {
	Area() float64
}

// структура Rectangle
type Rectangle struct {
	Width  float64 
	Height float64 
}

// метод Area для прямоугольника
// площадь прямоугольника : ширина * высота
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// структура Circle
type Circle struct {
	radius float64 
}

// метод Area для круга
// площадь круга : pi * r^2
func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

// функция принимает срез фигур shapes
func PrintAreas(shapes []Shape) {
	for _, shape := range shapes {
		fmt.Printf("Площадь фигуры: %.2f\n", shape.Area())
	}
}

func main() {
	var width, height float64 

	// запрос ширины прямоугольника
	fmt.Print("Введите ширину прямоугольника: ")
	_, err := fmt.Scan(&width)
	if err != nil {
		fmt.Println("Ошибка ввода ширины:", err)
		return
	}

	// запрос высоты прямоугольника
	fmt.Print("Введите высоту прямоугольника: ")
	_, err = fmt.Scan(&height)
	if err != nil {
		fmt.Println("Ошибка ввода высоты:", err)
		return
	}

	// прямоугольник с введенными высотой и шириной
	rect := Rectangle{
		Width:  width,
		Height: height,
	}

	var radius float64 

	// запрос радиуса круга
	fmt.Print("Введите радиус круга: ")
	_, err = fmt.Scan(&radius)
	if err != nil {
		fmt.Println("Ошибка ввода радиуса:", err)
		return
	}

	// круг с введенным радиусом 
	circle := Circle{
		radius: radius,
	}

	// срез для обоих приколов
	shapes := []Shape{rect, circle}

	// общий вывод
	PrintAreas(shapes)
}
