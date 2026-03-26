package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Account struct {
	Bank    string
	NoRek   string
	Balance int
}

type Transaction struct {
	ID     int
	From   string
	To     string
	Amount int
}

type TxResult struct {
	TxID    int
	Success bool
	Message string
}

var (
	accounts = map[string]*Account{
		"C001": {Bank: "CIMB", NoRek: "C001", Balance: 300000},
		"M002": {Bank: "MANDIRI", NoRek: "M002", Balance: 500000},
		"B003": {Bank: "BNI", NoRek: "B003", Balance: 400000},
		"BC04": {Bank: "BCA", NoRek: "BC04", Balance: 800000},
	}

	transactions = []Transaction{
		{ID: 1, From: "C001", To: "M002", Amount: 300000},
		{ID: 2, From: "C001", To: "B003", Amount: 600000},
		{ID: 3, From: "M002", To: "BC04", Amount: 200000},
		{ID: 4, From: "B003", To: "C001", Amount: 500000},
		{ID: 5, From: "BC04", To: "M002", Amount: 700000},
		{ID: 6, From: "C001", To: "M002", Amount: 400000},
		{ID: 7, From: "B003", To: "C001", Amount: 300000},
	}

	mu sync.Mutex
)

// Setiap worker mengambil transaksi dari channel jobs
func worker(id int, jobs <-chan Transaction, results chan<- TxResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for tx := range jobs {
		fmt.Printf("Worker-%d | Mengambil TX-%d\n", id, tx.ID)
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		processTransaction(ctx, tx, results)
		cancel()
	}
}

func processTransaction(ctx context.Context, tx Transaction, resultCh chan<- TxResult) {
	start := time.Now()
	now := start.Format("15:04:05")
	fmt.Printf("%s | Processing TX-%d | %s -> %s | Rp%d\n", now, tx.ID, tx.From, tx.To, tx.Amount)

	processingTime := time.Duration(1000+rand.Intn(4500)) * time.Millisecond

	select {
	case <-time.After(processingTime):

		mu.Lock()
		defer mu.Unlock()

		sender := accounts[tx.From]
		receiver := accounts[tx.To]
		now = time.Now().Format("15:04:05")

		if sender.Balance < tx.Amount {
			elapsed := time.Since(start).Seconds()
			msg := fmt.Sprintf("%s | TX-%d FAILED (%.1fs) | %s -> %s | Transfer Rp%d | Saldo %s hanya Rp%d",
				now, tx.ID, elapsed, tx.From, tx.To, tx.Amount, tx.From, sender.Balance)
			fmt.Println(msg)
			resultCh <- TxResult{TxID: tx.ID, Success: false, Message: msg}
			return
		}

		sender.Balance -= tx.Amount
		receiver.Balance += tx.Amount

		elapsed := time.Since(start).Seconds()
		fmt.Printf("%s | TX-%d SUCCESS (%.1fs)\n", now, tx.ID, elapsed)
		fmt.Printf("    Saldo %s sekarang Rp%d\n", sender.NoRek, sender.Balance)
		fmt.Printf("    Saldo %s sekarang Rp%d\n", receiver.NoRek, receiver.Balance)
		resultCh <- TxResult{TxID: tx.ID, Success: true}

	case <-ctx.Done():
		elapsed := time.Since(start).Seconds()
		now = time.Now().Format("15:04:05")
		fmt.Printf("%s | TX-%d TIMEOUT (%.1fs)\n", now, tx.ID, elapsed)
		resultCh <- TxResult{TxID: tx.ID, Success: false, Message: "timeout"}
	}
}

func main() {

	jobs := make(chan Transaction, len(transactions))
	results := make(chan TxResult, len(transactions))

	var wg sync.WaitGroup
	numWorkers := len(transactions)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for _, tx := range transactions {
		jobs <- tx
	}
	close(jobs)

	wg.Wait()
	close(results)

	for range results {
	}

	fmt.Println("\n===== FINAL BALANCE =====")
	order := []string{"C001", "M002", "B003", "BC04"}
	total := 0
	for _, noRek := range order {
		acc := accounts[noRek]
		fmt.Printf("%s (%s) : Rp%d\n", acc.Bank, acc.NoRek, acc.Balance)
		total += acc.Balance
	}
	fmt.Printf("\nTotal Saldo Semua Bank : Rp%d\n", total)
}
