package main

import (
	"fmt"
	"math"
)

type Circle struct {
	radius float64
}

// метод для вычисления площади круга
func (c Circle) area() float64 {
	// формула площади круга: pi * r^2
	return math.Pi * c.radius * c.radius
}

func main() {
	// создание экземпляра структуры Circle
	circle := Circle{radius: 3.0}

	// результат
	fmt.Printf("Площадь круга с радиусом %.2f равна %.2f\n", circle.radius, circle.area())
}
