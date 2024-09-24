package stringutils

// функция для переворота строки
func Reverse(s string) string {
	runes := []rune(s) // преобразование строки в срез рун для работы с символами
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i] // замена местами
	}
	return string(runes) // преобразование обратно в строку
}
