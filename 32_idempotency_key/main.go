package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func processTransaction(ctx context.Context, rdb *redis.Client, idempotencyKey string, amount int) error {
	// Dengan SetNX, command ini hanya berhasil meng-set value ke Redis JIKA key belum ada.
	// Jika key berhasil di set, artinya key tersebut belum pernah ada sebelumnya.
	// Jika return false, operasi duplicate / payment sedang/telah diproses.

	isFirstTime, err := rdb.SetNX(ctx, idempotencyKey, "started", 10*time.Minute).Result()
	if err != nil {
		return err
	}

	if !isFirstTime {
		fmt.Printf("[DITOLAK] Idempotency Key '%s' pernah/sedang di proses. Mencegah DEDUPLIKASI transfer!\n", idempotencyKey)
		return nil
	}

	fmt.Printf("[BERHASIL] Memproses transfer Rp %d \n", amount)
	return nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	ctx := context.Background()

	// Pastikan kita membersihkan redis untuk simulasi
	rdb.FlushDB(ctx)

	clientPaymentKey := "payment-req-39a0x-9f22p"

	fmt.Println("==> Percobaan Pertama ==")
	err := processTransaction(ctx, rdb, clientPaymentKey, 500000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n==> Percobaan Kedua (Misal Client Timeout dan men-klik tombol Pay lagi) ==")
	err = processTransaction(ctx, rdb, clientPaymentKey, 500000)
	if err != nil {
		log.Fatal(err)
	}
}
