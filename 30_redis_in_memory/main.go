package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Struct untuk merepresentasikan data User
type User struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
}

func main() {
	// Koneksi ke Docker Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	// 1. Data Asli
	ayash := User{
		Name:       "Ayash",
		Age:        21,
		Occupation: "Software Engineer",
	}

	fmt.Println("========= 1. MENYIMPAN DATA KE REDIS =========")
	// Convert Go Struct ke JSON String agar bisa disimpan ke Redis
	jsonData, err := json.Marshal(ayash)
	if err != nil {
		log.Fatalf("Gagal konversi ke JSON: %v", err)
	}

	cacheKey := "user_profile:ayash"
	// Simpan JSON string ke Redis
	err = rdb.Set(ctx, cacheKey, jsonData, 0).Err()
	if err != nil {
		log.Fatalf("Gagal menyimpan ke Redis: %v", err)
	}
	fmt.Printf("Data Go berhasil di-convert ke JSON dan masuk ke Redis!\n")
	fmt.Printf("Key   : %s\n", cacheKey)
	fmt.Printf("Value : %s\n", string(jsonData))

	fmt.Println("\n========= 2. MENGAMBIL DATA DARI REDIS =========")

	// Ambil Data dari Redis
	// (Seolah-olah database sedang down dan kita mengandalkan Redis)
	redisValue, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		fmt.Println("[CACHE MISS] Data tidak ditemukan di Redis")
		return
	} else if err != nil {
		log.Fatalf("Gagal mengambil dari Redis: %v", err)
	}

	// Tanda jelas bahwa data ini murni string mentah yang ditarik dari Redis cache
	fmt.Println("[CACHE HIT] Berhasil narik string mentah dari sistem Redis:")
	fmt.Println("->", redisValue)
	fmt.Println()

	// Convert kembali JSON String mentah dari Redis ke Go Struct asli
	var retrievedUser User
	err = json.Unmarshal([]byte(redisValue), &retrievedUser)
	if err != nil {
		log.Fatalf("Gagal parse JSON dari redis: %v", err)
	}

	fmt.Println("[BERHASIL DI-PARSE KEMBALI KE GO STRUCT]")
	fmt.Printf("Nama      : %s\n", retrievedUser.Name)
	fmt.Printf("Umur      : %d\n", retrievedUser.Age)
	fmt.Printf("Pekerjaan : %s\n", retrievedUser.Occupation)
}
