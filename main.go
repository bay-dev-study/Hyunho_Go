package main

import (
	"Hyunho_Go/dictionary"
	"log"
)

const PASS_ERROR = false
const CHECK_ERROR = true

func errCheck(err error, flag bool) {
	if err != nil && flag == CHECK_ERROR {
		log.Fatalln(err)
	}
}

func main() {
	dct := dictionary.Dictionary{}
	dct.Add("Hello", "Greeting")
	dct.Add("Apple", "Company")
	dct.Add("Ant", "Insect")
	dct.Print()

	err := dct.Add("Apple", "Fruit") // raise error
	errCheck(err, PASS_ERROR)

	dct.Update("Apple", "Fruit")
	dct.Print()

	err = dct.Update("Samsung", "Company") // raise error
	errCheck(err, PASS_ERROR)

	dct.Delete("Apple")

	_, err = dct.Search("Apple") // raise error
	errCheck(err, PASS_ERROR)

	dct.Print()
}
