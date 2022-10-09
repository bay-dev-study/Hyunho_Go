package banking

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	owner   string
	balance float64
}

func NewAccount(owner string, balance float64) *BankAccount {
	return &BankAccount{owner, balance}
}

func (bankAccount *BankAccount) Deposit(depositAmount float64) error {
	if depositAmount <= 0 {
		return errors.New("wrong deposit amount")
	}
	bankAccount.balance += depositAmount
	return nil
}

func (bankAccount *BankAccount) Withdraw(withdrawAmount float64) error {
	if withdrawAmount <= 0 || withdrawAmount >= bankAccount.balance {
		return errors.New("wrong withdraw amount")
	}
	bankAccount.balance -= withdrawAmount
	return nil
}

func (bankAccount *BankAccount) Balance() float64 {
	return bankAccount.balance
}

func (bankAccount *BankAccount) Owner() string {
	return bankAccount.owner
}

func (bankAccount *BankAccount) String() string {
	return fmt.Sprintf("Owner: %s, Balance: %f", bankAccount.owner, bankAccount.balance)
}
