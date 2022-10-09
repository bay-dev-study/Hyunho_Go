package main

import (
	"Hyunho_Go/banking"
	"fmt"
	"log"
)

func main() {
	bankAccount := banking.NewAccount("Hyunho", 100)
	fmt.Println(bankAccount.Owner(), bankAccount.Balance())
	fmt.Println(bankAccount)
	bankAccount.Deposit(10)
	fmt.Println(bankAccount)
	bankAccount.Withdraw(70)
	fmt.Println(bankAccount)
	err := bankAccount.Withdraw(50)
	if err != nil {
		log.Fatalln(err)
	}
}
