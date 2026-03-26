# 26. Synchronization — Channels, Select, Worker Pools, Context

Modul ini membahas mekanisme **synchronization** di Go untuk koordinasi antar goroutine.

## Channels

### Unbuffered Channel
```go
ch := make(chan string)       // synchronous
go func() { ch <- "hello" }() // blocks until receiver ready
msg := <-ch                    // blocks until sender sends
```

### Buffered Channel
```go
ch := make(chan int, 3)  // buffer size 3
ch <- 1                  // doesn't block (buffer not full)
ch <- 2
ch <- 3
// ch <- 4              // BLOCKS — buffer full!
```

### Channel Direction
```go
func send(ch chan<- string) { ch <- "data" }  // send only
func recv(ch <-chan string) { <-ch }          // receive only
```

### Range Over Channel
```go
ch := make(chan int)
go func() {
    for i := 0; i < 5; i++ { ch <- i }
    close(ch) // WAJIB close agar range berhenti
}()
for val := range ch { fmt.Println(val) }
```

## Select Statement

### Basic Select
```go
select {
case msg := <-ch1:
    fmt.Println(msg)
case msg := <-ch2:
    fmt.Println(msg)
}
```

### Timeout
```go
select {
case msg := <-ch:
    fmt.Println(msg)
case <-time.After(5 * time.Second):
    fmt.Println("timeout!")
}
```

### Non-blocking (default)
```go
select {
case msg := <-ch:
    fmt.Println(msg)
default:
    fmt.Println("no message yet")
}
```

## Worker Pool Pattern

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2
    }
}

jobs := make(chan int, 100)
results := make(chan int, 100)

// Start 3 workers
for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}
```

## Context

### WithCancel
```go
ctx, cancel := context.WithCancel(context.Background())
go worker(ctx)
cancel() // signal stop
```

### WithTimeout
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### WithDeadline
```go
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

### WithValue
```go
ctx := context.WithValue(ctx, "requestID", "abc-123")
reqID := ctx.Value("requestID")
```

## Key Takeaways

| Tool | Kegunaan |
|------|----------|
| Channel | Komunikasi antar goroutine |
| Buffered Channel | Async communication (sampai buffer full) |
| Select | Multiplex beberapa channel |
| time.After | Timeout pada channel operation |
| Worker Pool | Fixed goroutines + job queue |
| context.WithCancel | Manual cancellation |
| context.WithTimeout | Auto-cancel setelah durasi |
| context.WithDeadline | Auto-cancel pada waktu tertentu |
| context.WithValue | Pass request-scoped data |
