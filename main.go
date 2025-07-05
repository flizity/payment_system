package main

import (
	"errors"
	"fmt"
	"time"
)

type Acc struct {
	ID      int
	Name    string
	Balance float64
}

func (u *Acc) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("нельзя вносить отрицательную сумму")
	}
	u.Balance += amount
	return nil
}

func (u *Acc) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("сумма должна быть положительной")
	}
	if u.Balance < amount {
		return errors.New("недостаточно средств")
	}
	u.Balance -= amount
	return nil
}

type Transaction struct {
	FromID int
	ToID   int
	Amount float64
	Time   time.Time
	Status string
}

type PaymentSystem struct {
	Users            map[int]*Acc
	TransactionQueue []Transaction
}

func (ps *PaymentSystem) AddUser(user *Acc) {
	if ps.Users == nil {
		ps.Users = make(map[int]*Acc)
	}
	ps.Users[user.ID] = user
}

func (ps *PaymentSystem) AddTransaction(fromID, toID int, amount float64) {
	ps.TransactionQueue = append(ps.TransactionQueue, Transaction{
		FromID: fromID,
		ToID:   toID,
		Amount: amount,
		Time:   time.Now(),
		Status: "ожидает обработки",
	})
}

func (ps *PaymentSystem) ProcessTransactions() {
	for i := range ps.TransactionQueue {
		t := &ps.TransactionQueue[i]
		if t.Status != "ожидает обработки" {
			continue
		}

		fromUser, existsFrom := ps.Users[t.FromID]
		toUser, existsTo := ps.Users[t.ToID]

		if !existsFrom || !existsTo {
			t.Status = "ошибка: пользователь не найден"
			continue
		}

		if err := fromUser.Withdraw(t.Amount); err != nil {
			t.Status = "ошибка: " + err.Error()
			continue
		}

		if err := toUser.Deposit(t.Amount); err != nil {
			fromUser.Deposit(t.Amount) // Откат транзакции
			t.Status = "ошибка: " + err.Error()
			continue
		}

		t.Status = "успешно выполнена"
	}
}

func main() {
	u1 := &Acc{14773, "Alex", 2000}
	u2 := &Acc{23432, "Max", 15000}

	ps := PaymentSystem{}
	ps.AddUser(u1)
	ps.AddUser(u2)

	ps.AddTransaction(u1.ID, u2.ID, 100.50)
	ps.AddTransaction(u2.ID, u1.ID, 5000)

	ps.ProcessTransactions()

	fmt.Println("История транзакций:")
	for _, t := range ps.TransactionQueue {
		fmt.Printf("[%s] %d → %d: $%.2f (%s)\n",
			t.Time.Format("15:04:05"),
			t.FromID,
			t.ToID,
			t.Amount,
			t.Status)
	}

	fmt.Printf("\nБаланс %s: %.2f\n", u1.Name, u1.Balance)
	fmt.Printf("Баланс %s: %.2f\n", u2.Name, u2.Balance)
}
