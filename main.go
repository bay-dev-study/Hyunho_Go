package main

import "fmt"

func addNumbers(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}
func main() {
	fmt.Println(addNumbers(1, 2, 3, 4, 5))
}
