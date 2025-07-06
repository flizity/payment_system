package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Acc struct {
	ID      int
	Name    string
	Balance float64
	mu      sync.Mutex
}

func (u *Acc) Deposit(amount float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if amount <= 0 {
		return errors.New("нельзя вносить отрицательную сумму")
	}
	u.Balance += amount
	return nil
}

func (u *Acc) Withdraw(amount float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
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

func (ps *PaymentSystem) ProcessTransaction(t *Transaction) error {
	if t.Status != "ожидает обработки" {
		return nil
	}

	fromUser, existsFrom := ps.Users[t.FromID]
	toUser, existsTo := ps.Users[t.ToID]

	if !existsFrom || !existsTo {
		t.Status = "ошибка: пользователь не найден"
		return fmt.Errorf("ошибка: пользователь не найден")
	}

	if err := fromUser.Withdraw(t.Amount); err != nil {
		t.Status = "ошибка: " + err.Error()
		return err
	}

	if err := toUser.Deposit(t.Amount); err != nil {
		fromUser.Deposit(t.Amount)
		t.Status = "ошибка: " + err.Error()
		return err
	}

	t.Status = "успешно выполнена"
	return nil
}

func Worker(ps *PaymentSystem, ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		err := ps.ProcessTransaction(&t)
		if err != nil {
			fmt.Println("Ошибка обработки транзакции:", err)
		}
	}
}

func (ps *PaymentSystem) ProcessTransactions() {
	var wg sync.WaitGroup
	transactionChan := make(chan Transaction, len(ps.TransactionQueue))

	workerCount := 5
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go Worker(ps, transactionChan, &wg)
	}

	for _, t := range ps.TransactionQueue {
		transactionChan <- t
	}
	close(transactionChan)
	wg.Wait()
}

func main() {
	u1 := &Acc{ID: 14773, Name: "Alex", Balance: 2000}
	u2 := &Acc{ID: 23432, Name: "Max", Balance: 15000}

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
