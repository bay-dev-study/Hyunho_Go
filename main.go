package main

import "fmt"

func canIDrink(age int) bool {
	if koreanAge := age + 2; koreanAge > 18 {
		return true
	}
	return false
}

func main() {
	fmt.Println(canIDrink(16))
}
