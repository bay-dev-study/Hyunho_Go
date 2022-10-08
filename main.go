package main

import (
	"fmt"
	"strings"
)

func lengthAndUpper(name string) (length int, upperLetters string) {
	defer fmt.Println("I'm Done")
	length = len(name)
	upperLetters = strings.ToUpper(name)
	return
}

func main() {
	length, upperLetters := lengthAndUpper("Hello")
	fmt.Println(length, upperLetters)
}
