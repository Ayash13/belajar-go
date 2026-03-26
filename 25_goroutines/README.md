# 25. Go Routines

Modul ini membahas **goroutine** — unit eksekusi ringan yang menjadi fondasi concurrency di Go.

## Apa itu Goroutine?

Goroutine adalah **lightweight thread** yang dikelola oleh Go runtime, bukan OS. Sangat murah untuk dibuat (stack ~2-8 KB).

```go
// Buat goroutine dengan keyword 'go'
go func() {
    fmt.Println("Hello from goroutine!")
}()
```

## Yang Didemonstrasikan

### 1. Basic Goroutine
```go
go func() {
    fmt.Println("running in background")
}()
```

### 2. sync.WaitGroup
```go
var wg sync.WaitGroup
wg.Add(1)          // tambah counter
go func() {
    defer wg.Done() // kurangi counter saat selesai
    doWork()
}()
wg.Wait()           // tunggu sampai counter = 0
```

### 3. Closure Pitfall
```go
// ❌ WRONG — semua goroutine lihat value terakhir
for i := 0; i < 5; i++ {
    go func() { fmt.Println(i) }() // i = 5 semua!
}

// ✅ RIGHT — pass sebagai parameter
for i := 0; i < 5; i++ {
    go func(val int) { fmt.Println(val) }(i)
}
```

### 4. Race Condition & Mutex
```go
var mu sync.Mutex
counter := 0

go func() {
    mu.Lock()
    counter++
    mu.Unlock()
}()
```

### 5. Atomic Operations
```go
var counter int64
atomic.AddInt64(&counter, 1) // lock-free increment
```

## Goroutine vs OS Thread

| Feature | Goroutine | OS Thread |
|---------|-----------|-----------|
| Stack Size | ~2-8 KB | ~1-8 MB |
| Creation Time | ~μs | ~ms |
| Managed By | Go Runtime | OS Kernel |
| Context Switch | Sangat cepat | Relatif mahal |
| Max Count | Jutaan | Ribuan |

## Key Takeaways

- `go func(){}()` untuk launch goroutine
- `sync.WaitGroup` untuk menunggu goroutine selesai
- `sync.Mutex` untuk proteksi shared data
- `sync/atomic` untuk operasi atomic tanpa lock
- **Selalu pass loop variable** sebagai parameter ke goroutine
