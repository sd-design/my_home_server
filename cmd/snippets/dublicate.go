package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	// Создадим пустую map
	m := make(map[int]string)
	final := make(map[int]string)
	i := 0
	for _, tmp := range input {
		m[i] += tmp
		if _, ok := m[i]; !ok {
			final = append(final, i)

		}
		i++
	}
	fmt.Println(final)
	return nil
}

// как можно заметить, алгоритм пройдётся по массиву всего один раз
// если бы мы искали подходящее значение каждый раз через перебор массива, то пришлось бы сделать гораздо больше вычислений

func main() {
	input := []string{
		"cat",
		"dog",
		"bird",
		"dog",
		"parrot",
		"cat",
	}
	RemoveDuplicates(input)
}
