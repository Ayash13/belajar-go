package goroutines

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ══════════════════════════════════════════════════
// Go Routines — Lightweight Concurrent Execution
// ══════════════════════════════════════════════════

func Run() {
	fmt.Println("=== Go Routines ===")

	// ── 1. Basic Goroutine ──
	fmt.Println("\n── 1. Basic Goroutine ──")
	fmt.Println("   Goroutine = lightweight thread managed by Go runtime")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("   Hello from goroutine!")
	}()
	wg.Wait()
	fmt.Println("   Main goroutine continues after child finishes")

	// ── 2. Multiple Goroutines ──
	fmt.Println("\n── 2. Multiple Goroutines ──")
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id*20) * time.Millisecond)
			fmt.Printf("   Goroutine %d selesai\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("   Semua goroutine selesai")

	// ── 3. WaitGroup ──
	fmt.Println("\n── 3. sync.WaitGroup ──")
	fmt.Println("   WaitGroup = counter untuk menunggu goroutines selesai")

	var wg2 sync.WaitGroup
	jobs := []string{"download", "process", "upload"}
	for _, job := range jobs {
		wg2.Add(1) // increment counter
		go func(j string) {
			defer wg2.Done() // decrement counter saat selesai
			fmt.Printf("   [START] %s\n", j)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("   [DONE]  %s\n", j)
		}(job)
	}
	wg2.Wait() // block sampai counter = 0
	fmt.Println("   All jobs completed!")

	// ── 4. Anonymous Goroutine ──
	fmt.Println("\n── 4. Anonymous Goroutine ──")
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("   Ini anonymous goroutine (fungsi tanpa nama)")
	}()
	wg.Wait()

	// ── 5. Goroutine with Return Value (via channel) ──
	fmt.Println("\n── 5. Goroutine Return Value ──")
	fmt.Println("   Goroutine tidak bisa return value langsung")
	fmt.Println("   Gunakan channel untuk mengirim hasil")

	result := make(chan int)
	go func() {
		sum := 0
		for i := 1; i <= 100; i++ {
			sum += i
		}
		result <- sum // kirim hasil via channel
	}()
	total := <-result // terima hasil
	fmt.Printf("   Sum 1-100 = %d\n", total)

	// ── 6. Goroutine Closure Pitfall ──
	fmt.Println("\n── 6. Closure Pitfall ──")
	fmt.Println("   ❌ WRONG: variable di-capture by reference")
	// Demonstrasi yang BENAR
	fmt.Println("   ✅ RIGHT: pass variable sebagai parameter")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			fmt.Printf("   val = %d\n", val)
		}(i) // pass i sebagai argument
	}
	wg.Wait()

	// ── 7. Goroutine vs Thread ──
	fmt.Println("\n── 7. Goroutine vs OS Thread ──")
	fmt.Println("   ┌───────────────┬──────────────┬──────────────┐")
	fmt.Println("   │ Feature       │ Goroutine    │ OS Thread    │")
	fmt.Println("   ├───────────────┼──────────────┼──────────────┤")
	fmt.Println("   │ Stack Size    │ ~2-8 KB      │ ~1-8 MB      │")
	fmt.Println("   │ Creation      │ ~μs          │ ~ms          │")
	fmt.Println("   │ Managed By    │ Go Runtime   │ OS Kernel    │")
	fmt.Println("   │ Switching     │ Sangat cepat │ Relatif mahal│")
	fmt.Println("   │ Jumlah        │ Jutaan       │ Ribuan       │")
	fmt.Println("   └───────────────┴──────────────┴──────────────┘")

	// ── 8. Spawning Many Goroutines ──
	fmt.Println("\n── 8. Spawning 10,000 Goroutines ──")
	var counter int64
	var wg3 sync.WaitGroup
	start := time.Now()
	for i := 0; i < 10_000; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg3.Wait()
	elapsed := time.Since(start)
	fmt.Printf("   Spawned & completed: %d goroutines\n", counter)
	fmt.Printf("   Time: %v\n", elapsed.Round(time.Microsecond))

	// ── 9. Race Condition Demo ──
	fmt.Println("\n── 9. Race Condition & Mutex ──")
	fmt.Println("   Tanpa synchronization, data bisa corrupt")

	// Safe with Mutex
	var mu sync.Mutex
	safeCounter := 0
	var wg4 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg4.Add(1)
		go func() {
			defer wg4.Done()
			mu.Lock()
			safeCounter++
			mu.Unlock()
		}()
	}
	wg4.Wait()
	fmt.Printf("   Counter with Mutex: %d (correct!)\n", safeCounter)

	// Safe with Atomic
	var atomicCounter int64
	var wg5 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg5.Add(1)
		go func() {
			defer wg5.Done()
			atomic.AddInt64(&atomicCounter, 1)
		}()
	}
	wg5.Wait()
	fmt.Printf("   Counter with Atomic: %d (correct!)\n", atomicCounter)

	// ── Summary ──
	fmt.Println("\n=== Summary ===")
	fmt.Println("  go func() { ... }()  → launch goroutine")
	fmt.Println("  sync.WaitGroup       → wait for goroutines to finish")
	fmt.Println("  sync.Mutex           → protect shared data")
	fmt.Println("  sync/atomic          → lock-free atomic operations")
	fmt.Println("  Always pass loop vars as params to avoid closure bugs")
}
