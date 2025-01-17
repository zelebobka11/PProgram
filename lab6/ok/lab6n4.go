package main

import (
	"fmt"
	"sync"
)

// общий счетчик
var counter,bebr int

// мьютекс для синхронизации доступа к общему ресурсу
var mutex sync.Mutex

// кол-во горутин
const goroutines = 5

// функция для увеличения счётчика
func increment(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10000; i++ {
		mutex.Lock()   // блокируем доступ к счетчику               
		counter++      // увеличиваем общий счётчик
		mutex.Unlock() // разблокируем доступ к счётчику            
	}
	
	for i := 0; i < 10000; i++ {
		
		bebr++      // увеличиваем общий счётчик
		
	}

}

func main() {
	// инициализация WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// запуск горутин
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go increment(&wg)
	}

	// ожидание завершения всех горутин
	wg.Wait()

	// вывод конечного значения счётчика
	fmt.Printf("mu: %d\n", counter)
	fmt.Printf("no mu: %d\n", bebr)
}
