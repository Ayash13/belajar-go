# 24. Concurrency vs Parallelism

Modul ini menjelaskan perbedaan fundamental antara **concurrency** dan **parallelism** di Go.

## Konsep Dasar

### Concurrency (Konkurensi)
- **Structure** — cara meng-*organize* program untuk menangani banyak task
- Tasks bisa START, RUN, dan COMPLETE dalam waktu yang **overlapping**
- Bisa dilakukan bahkan di **1 CPU core**
- Contoh: 1 chef memasak beberapa hidangan (bergantian mengerjakan)

### Parallelism (Paralelisme)
- **Execution** — tasks benar-benar berjalan **bersamaan** di waktu yang sama
- Membutuhkan **multiple CPU cores**
- Contoh: beberapa chef, masing-masing memasak 1 hidangan secara simultan

```
Concurrency ≠ Parallelism

Concurrency = dealing with lots of things at once (design)
Parallelism  = doing lots of things at once (execution)
```

## Yang Didemonstrasikan

### 1. Runtime Info
```go
runtime.NumCPU()        // jumlah CPU core
runtime.GOMAXPROCS(0)   // jumlah OS thread untuk goroutines
runtime.NumGoroutine()  // jumlah goroutine aktif
```

### 2. Sequential vs Concurrent
```go
// Sequential: A → B → C (total = 200+150+100 = 450ms)
simulateTask("A", 200ms)
simulateTask("B", 150ms)
simulateTask("C", 100ms)

// Concurrent: A, B, C run together (total ≈ 200ms)
go simulateTask("A", 200ms)
go simulateTask("B", 150ms)
go simulateTask("C", 100ms)
```

### 3. GOMAXPROCS
```go
runtime.GOMAXPROCS(1)            // 1 thread → concurrent only
runtime.GOMAXPROCS(runtime.NumCPU()) // all cores → concurrent + parallel
```

## Key Takeaways

| Aspect | Concurrency | Parallelism |
|--------|------------|-------------|
| Definisi | Menangani banyak task | Menjalankan banyak task bersamaan |
| CPU Core | 1 core cukup | Butuh multiple cores |
| Go Tool | `go` keyword (goroutine) | `GOMAXPROCS` |
| Analogi | 1 chef, banyak hidangan | Banyak chef, banyak hidangan |
