package concurrency

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// ══════════════════════════════════════════════════
// Concurrency vs Parallelism
// ══════════════════════════════════════════════════

// simulateTask simulates work with a given duration.
func simulateTask(name string, duration time.Duration) {
	fmt.Printf("   [START] %s\n", name)
	time.Sleep(duration)
	fmt.Printf("   [DONE]  %s (took %v)\n", name, duration)
}

// sequentialWork runs tasks one after another.
func sequentialWork() time.Duration {
	start := time.Now()
	simulateTask("Task A", 200*time.Millisecond)
	simulateTask("Task B", 150*time.Millisecond)
	simulateTask("Task C", 100*time.Millisecond)
	return time.Since(start)
}

// concurrentWork runs tasks concurrently using goroutines.
func concurrentWork() time.Duration {
	start := time.Now()
	var wg sync.WaitGroup

	tasks := []struct {
		name     string
		duration time.Duration
	}{
		{"Task A", 200 * time.Millisecond},
		{"Task B", 150 * time.Millisecond},
		{"Task C", 100 * time.Millisecond},
	}

	for _, t := range tasks {
		wg.Add(1)
		go func(name string, dur time.Duration) {
			defer wg.Done()
			simulateTask(name, dur)
		}(t.name, t.duration)
	}

	wg.Wait()
	return time.Since(start)
}

func Run() {
	fmt.Println("=== Concurrency vs Parallelism ===")

	// ── 1. CPU Info ──
	fmt.Println("\n── 1. Runtime Info ──")
	fmt.Printf("   NumCPU:       %d\n", runtime.NumCPU())
	fmt.Printf("   GOMAXPROCS:   %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("   NumGoroutine: %d\n", runtime.NumGoroutine())

	// ── 2. Sequential Execution ──
	fmt.Println("\n── 2. Sequential Execution ──")
	fmt.Println("   Tasks run one after another (blocking)")
	seqDur := sequentialWork()
	fmt.Printf("   Total time: %v\n", seqDur.Round(time.Millisecond))

	// ── 3. Concurrent Execution ──
	fmt.Println("\n── 3. Concurrent Execution ──")
	fmt.Println("   Tasks run concurrently (non-blocking)")
	conDur := concurrentWork()
	fmt.Printf("   Total time: %v\n", conDur.Round(time.Millisecond))

	// ── 4. Speedup ──
	fmt.Println("\n── 4. Speedup Comparison ──")
	fmt.Printf("   Sequential:  %v\n", seqDur.Round(time.Millisecond))
	fmt.Printf("   Concurrent:  %v\n", conDur.Round(time.Millisecond))
	speedup := float64(seqDur) / float64(conDur)
	fmt.Printf("   Speedup:     %.2fx faster\n", speedup)

	// ── 5. Concurrency vs Parallelism Explained ──
	fmt.Println("\n── 5. Key Concepts ──")
	fmt.Println("   Concurrency:")
	fmt.Println("     → Dealing with multiple things at once (structure)")
	fmt.Println("     → Tasks can START, RUN, and COMPLETE in overlapping time")
	fmt.Println("     → Possible even on a single CPU core")
	fmt.Println("")
	fmt.Println("   Parallelism:")
	fmt.Println("     → Doing multiple things at once (execution)")
	fmt.Println("     → Tasks LITERALLY run at the same time on different cores")
	fmt.Println("     → Requires multiple CPU cores")
	fmt.Println("")
	fmt.Println("   Analogy:")
	fmt.Println("     Concurrency = 1 chef, multiple dishes (switching between tasks)")
	fmt.Println("     Parallelism = multiple chefs, each cooking a dish simultaneously")

	// ── 6. GOMAXPROCS Demo ──
	fmt.Println("\n── 6. GOMAXPROCS Demo ──")

	// Single core
	runtime.GOMAXPROCS(1)
	fmt.Printf("   GOMAXPROCS(1): single OS thread → concurrent but NOT parallel\n")
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sum := 0
			for j := 0; j < 1_000_000; j++ {
				sum += j
			}
			_ = sum
		}(i)
	}
	wg.Wait()
	single := time.Since(start)

	// All cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("   GOMAXPROCS(%d): all cores → concurrent AND parallel\n", runtime.NumCPU())
	start = time.Now()
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sum := 0
			for j := 0; j < 1_000_000; j++ {
				sum += j
			}
			_ = sum
		}(i)
	}
	wg.Wait()
	multi := time.Since(start)

	fmt.Printf("   1 core:    %v\n", single.Round(time.Microsecond))
	fmt.Printf("   %d cores:   %v\n", runtime.NumCPU(), multi.Round(time.Microsecond))

	// ── Summary ──
	fmt.Println("\n=== Summary ===")
	fmt.Println("  Concurrency = STRUCTURE (managing multiple tasks)")
	fmt.Println("  Parallelism = EXECUTION (running tasks simultaneously)")
	fmt.Println("  Go makes concurrency easy with goroutines")
	fmt.Println("  GOMAXPROCS controls how many OS threads run goroutines")
}
