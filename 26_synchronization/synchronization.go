package synchronization

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ══════════════════════════════════════════════════
// Synchronization — Channels, Select, Worker Pools, Context
// ══════════════════════════════════════════════════

func Run() {
	fmt.Println("=== Synchronization: Channels, Select, Worker Pools, Context ===")

	// ── 1. Unbuffered Channel ──
	fmt.Println("\n── 1. Unbuffered Channel ──")
	fmt.Println("   Sender blocks until receiver is ready (synchronous)")

	ch := make(chan string)
	go func() {
		ch <- "hello from goroutine"
	}()
	msg := <-ch
	fmt.Printf("   Received: %s\n", msg)

	// ── 2. Buffered Channel ──
	fmt.Println("\n── 2. Buffered Channel ──")
	fmt.Println("   Buffered = sender doesn't block until buffer is full")

	buffered := make(chan int, 3)
	buffered <- 10
	buffered <- 20
	buffered <- 30
	fmt.Printf("   Buffer len: %d, cap: %d\n", len(buffered), cap(buffered))
	fmt.Printf("   Values: %d, %d, %d\n", <-buffered, <-buffered, <-buffered)

	// ── 3. Channel Direction ──
	fmt.Println("\n── 3. Channel Direction ──")
	fmt.Println("   chan<- = send only, <-chan = receive only")

	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "message")
	pong(pings, pongs)
	fmt.Printf("   Pong received: %s\n", <-pongs)

	// ── 4. Range Over Channel ──
	fmt.Println("\n── 4. Range Over Channel ──")
	nums := make(chan int)
	go func() {
		for i := 1; i <= 5; i++ {
			nums <- i
		}
		close(nums) // MUST close to stop range
	}()
	fmt.Print("   Values: ")
	for n := range nums {
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// ── 5. Select Statement ──
	fmt.Println("\n── 5. Select Statement ──")
	fmt.Println("   Select = switch for channels, picks whichever is ready first")

	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "from ch1"
	}()
	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- "from ch2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("   Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("   Received: %s\n", msg2)
		}
	}

	// ── 6. Select with Timeout ──
	fmt.Println("\n── 6. Select with Timeout ──")
	slowCh := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		slowCh <- "slow response"
	}()

	select {
	case msg := <-slowCh:
		fmt.Printf("   Got: %s\n", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("   Timeout! No response within 100ms")
	}

	// ── 7. Select with Default (Non-blocking) ──
	fmt.Println("\n── 7. Select with Default ──")
	emptyCh := make(chan string)
	select {
	case msg := <-emptyCh:
		fmt.Printf("   Got: %s\n", msg)
	default:
		fmt.Println("   No message available (non-blocking check)")
	}

	// ── 8. Done Channel Pattern ──
	fmt.Println("\n── 8. Done Channel Pattern ──")
	done := make(chan bool)
	go func() {
		fmt.Println("   Worker: processing...")
		time.Sleep(50 * time.Millisecond)
		fmt.Println("   Worker: done!")
		done <- true
	}()
	<-done
	fmt.Println("   Main: worker finished")

	// ── 9. Fan-Out / Fan-In ──
	fmt.Println("\n── 9. Fan-Out / Fan-In ──")
	fmt.Println("   Fan-Out: satu producer, banyak consumer")
	fmt.Println("   Fan-In:  banyak producer, satu consumer")

	jobs := make(chan int, 5)
	results := make(chan string, 5)

	// Fan-Out: 3 workers consuming from 1 channel
	for w := 1; w <= 3; w++ {
		go fanOutWorker(w, jobs, results)
	}

	// Send jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// Fan-In: collect all results
	for i := 0; i < 5; i++ {
		fmt.Printf("   %s\n", <-results)
	}

	// ── 10. Worker Pool ──
	fmt.Println("\n── 10. Worker Pool Pattern ──")
	fmt.Println("   Fixed number of workers process jobs from a queue")

	const numWorkers = 3
	const numJobs = 8

	jobsCh := make(chan int, numJobs)
	resultsCh := make(chan WorkerResult, numJobs)

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go poolWorker(w, jobsCh, resultsCh, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobsCh <- j
	}
	close(jobsCh)

	// Wait and close results
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results
	for r := range resultsCh {
		fmt.Printf("   Worker %d processed Job %d → result: %d\n", r.WorkerID, r.JobID, r.Result)
	}

	// ── 11. Context — Cancellation ──
	fmt.Println("\n── 11. Context — Cancellation ──")
	fmt.Println("   context.WithCancel: manually cancel goroutines")

	ctx, cancel := context.WithCancel(context.Background())
	go contextWorker(ctx, "cancel-worker")
	time.Sleep(120 * time.Millisecond)
	cancel() // signal cancellation
	time.Sleep(20 * time.Millisecond)

	// ── 12. Context — Timeout ──
	fmt.Println("\n── 12. Context — Timeout ──")
	fmt.Println("   context.WithTimeout: auto-cancel after duration")

	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	go contextWorker(ctx2, "timeout-worker")
	time.Sleep(150 * time.Millisecond)

	// ── 13. Context — Deadline ──
	fmt.Println("\n── 13. Context — Deadline ──")
	fmt.Println("   context.WithDeadline: cancel at specific time")

	deadline := time.Now().Add(80 * time.Millisecond)
	ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
	defer cancel3()

	go contextWorker(ctx3, "deadline-worker")
	time.Sleep(120 * time.Millisecond)

	// ── 14. Context with Value ──
	fmt.Println("\n── 14. Context with Value ──")
	fmt.Println("   Pass request-scoped values (e.g., request ID)")

	ctx4 := context.WithValue(context.Background(), contextKey("requestID"), "req-abc-123")
	processRequest(ctx4)

	// ── Summary ──
	fmt.Println("\n=== Summary ===")
	fmt.Println("  Channels:")
	fmt.Println("    make(chan T)     → unbuffered (synchronous)")
	fmt.Println("    make(chan T, n)  → buffered (async until full)")
	fmt.Println("    close(ch)       → signal no more values")
	fmt.Println("    range ch        → receive until closed")
	fmt.Println("")
	fmt.Println("  Select:")
	fmt.Println("    select { case <-ch: } → multiplex channels")
	fmt.Println("    time.After(d)         → timeout")
	fmt.Println("    default               → non-blocking")
	fmt.Println("")
	fmt.Println("  Worker Pool:")
	fmt.Println("    Fixed goroutines + job channel + result channel")
	fmt.Println("")
	fmt.Println("  Context:")
	fmt.Println("    WithCancel   → manual cancellation")
	fmt.Println("    WithTimeout  → auto-cancel after duration")
	fmt.Println("    WithDeadline → auto-cancel at specific time")
	fmt.Println("    WithValue    → pass request-scoped data")
}

// ── Helper Functions ──

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func fanOutWorker(id int, jobs <-chan int, results chan<- string) {
	for j := range jobs {
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		results <- fmt.Sprintf("Worker %d finished job %d", id, j)
	}
}

type WorkerResult struct {
	WorkerID int
	JobID    int
	Result   int
}

func poolWorker(id int, jobs <-chan int, results chan<- WorkerResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		time.Sleep(time.Duration(30+rand.Intn(50)) * time.Millisecond)
		results <- WorkerResult{
			WorkerID: id,
			JobID:    j,
			Result:   j * j,
		}
	}
}

func contextWorker(ctx context.Context, name string) {
	for i := 1; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("   %s stopped: %v\n", name, ctx.Err())
			return
		default:
			time.Sleep(40 * time.Millisecond)
			fmt.Printf("   %s: tick %d\n", name, i)
		}
	}
}

type contextKey string

func processRequest(ctx context.Context) {
	reqID := ctx.Value(contextKey("requestID"))
	fmt.Printf("   Processing request: %v\n", reqID)
}
