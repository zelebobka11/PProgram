package main

import (
	"fmt"
	"sync"
)

// структура для запроса
type CalcRequest struct {
	Operation string      // операция (+, -, *, /)
	Operand1  float64     // первый операнд
	Operand2  float64     // ЭТО ВТОРОЙ!!!
	Result    chan float64 // канал для передачи результата
}

// сервер калькулятора, который обрабатывает запросы
func calculator(requests chan CalcRequest, wg *sync.WaitGroup) {
	defer wg.Done() // уменьшаем счетчик при завершении
	for req := range requests {
		switch req.Operation {
		case "+":
			req.Result <- req.Operand1 + req.Operand2
		case "-":
			req.Result <- req.Operand1 - req.Operand2
		case "*":
			req.Result <- req.Operand1 * req.Operand2
		case "/":
			if req.Operand2 != 0 {
				req.Result <- req.Operand1 / req.Operand2
			} else {
				fmt.Println("Ошибка: Деление на ноль")
				req.Result <- 0 // возврат 0 при делении на ноль
			}
		default:
			fmt.Println("Неизвестная операция:", req.Operation)
			req.Result <- 0 // возврат 0 при неизвестной операции
		}
	}
}

func main() {
	// канал для передачи запросов к калькулятору
	requests := make(chan CalcRequest)

	// группа ожидания для синхронизации горутин
	var wg sync.WaitGroup

	// запуск сервера калькулятора
	wg.Add(1)
	go calculator(requests, &wg)

	// создание и отправка нескольких запросов
	for _, req := range []CalcRequest{
		{Operation: "+", Operand1: 10, Operand2: 5, Result: make(chan float64)},
		{Operation: "-", Operand1: 10, Operand2: 5, Result: make(chan float64)},
		{Operation: "*", Operand1: 10, Operand2: 5, Result: make(chan float64)},
		{Operation: "/", Operand1: 10, Operand2: 5, Result: make(chan float64)},
		{Operation: "/", Operand1: 10, Operand2: 0, Result: make(chan float64)}, // 0?!
	} {
		// отправляем запрос в канал
		requests <- req

		// читаем результат из канала и выводим
		go func(res chan float64) {
			fmt.Printf("Результат: %.2f\n", <-res)
		}(req.Result)
	}

	// закрываем канал ПОСЛЕ отправки всех запросов!!! 
	close(requests)

	// ожидание завершения работы всех горутин
	wg.Wait()
}
