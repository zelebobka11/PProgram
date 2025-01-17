package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type Task struct {
	Line string
}

func reverseString(s string) string {
	runes := []rune(s) // преобразуем строку в массив рун
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func worker(tasks <-chan Task, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		reversed := reverseString(task.Line) // реверсируем строку
		results <- reversed                   // отправляем результат в канал
	}
}

func main() {
    // чтение количества воркеров от пользователя
    var numWorkers int
    fmt.Print("Введите количество воркеров: ")
    fmt.Scan(&numWorkers)

    // инициализация каналов
    tasks := make(chan Task)
    results := make(chan string)

    var wg sync.WaitGroup

    // запуск воркеров
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(tasks, results, &wg)
    }

    // чтение строк из файла
    go func() {
        file, err := os.Open("input.txt") // откройте файла
        if err != nil {
            fmt.Println("Ошибка при открытии файла:", err)
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            line := scanner.Text()
            tasks <- Task{Line: line} // отправка задачи в канал
        }
        close(tasks) // закрываем канал задач
    }()

    // воркеры напишут результаты в этот канал
    go func() {
        wg.Wait()        // ожидаем пока все воркеры завершат работу
        close(results)   // закрываем канал результатов
    }()

    fmt.Println("Результаты реверсирования:")
    for result := range results {
        fmt.Println(result) // вывод результата в консоль
    }
}
