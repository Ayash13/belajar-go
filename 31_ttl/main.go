package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	ctx := context.Background()

	key := "OTP:user123"
	value := "739210"
	ttl := 2 * time.Second

	// Set data menggunakan Redis dengan TTL 2 detik
	err := rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Fatalf("Gagal set: %v", err)
	}

	fmt.Printf("\nBerhasil nyimpan OTP: %s=%s dengan TTL 2 detik\n", key, value)

	val, _ := rdb.Get(ctx, key).Result()
	fmt.Println("Get Langsung:", val) // Akan mereturn OTP

	fmt.Println("\nMenunggu 3 Detik...........")
	time.Sleep(3 * time.Second)

	// Harusnya OTP sudah expired dan tidak ketemu di cache
	valAfter, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("Waktu HABIS!, data otomatis dihapus oleh sistem Redis.")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Get Setelah 3 detik: %s\n", valAfter)
	}
}
