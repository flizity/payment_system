package main

import (
	"errors"
	"fmt"
)

type Acc struct {
	ID      int
	Name    string
	Balance float32
}

type Transaction struct {
	ToID int
}

func (u *Acc) Deposit(amount float32) error {
	if amount <= 0 {
		return errors.New("нельзя вносить отрицательную сумму")
	}
	u.Balance += amount
	return nil
}

func (u *Acc) Withdraw(amount float32) error {
	if amount <= 0 {
		return errors.New("сумма должна быть положительной")
	}
	if u.Balance < amount {
		return errors.New("недостаточно средств")
	}
	u.Balance -= amount
	return nil
}

func main() {
	User1 := Acc{ID: 1, Name: "Alex", Balance: 2000}
	User2 := Acc{ID: 2, Name: "Max", Balance: 15000}

	if err := User1.Deposit(400); err != nil {
		fmt.Println("Ошибка у Alex:", err)
	}
	if err := User2.Deposit(-20000); err != nil {
		fmt.Println("Ошибка у Max:", err)
	}

	fmt.Println(User1)
	fmt.Println(User2)
}
