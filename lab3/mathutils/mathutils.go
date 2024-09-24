package mathutils

// функция для вычисления факториала
func Factorial(n int) int {
	if n < 0 {
		return -1 // возвращаем -1 для отрицательных чисел
	}
	if n == 0 {
		return 1 
	}
	factorial := 1
	for i := 1; i <= n; i++ {
		factorial *= i
	}
	return factorial
}
