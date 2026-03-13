# 3. Functions

## Basic Function

```go
func greet(name string) string {
    return "Hello, " + name + "!"
}
```

Format: `func namaFunction(parameter tipe) returnType { ... }`

## Multiple Return Values

Go bisa return **lebih dari satu nilai**. Ini sangat umum dipakai, terutama untuk return `(result, error)`.

```go
func divide(a, b float64) (float64, string) {
    if b == 0 {
        return 0, "cannot divide by zero"
    }
    return a / b, ""
}

// penggunaan:
result, errMsg := divide(10, 3)
```

## Named Return Values

Return value bisa dikasih nama. Cukup `return` tanpa argumen.

```go
func rectangle(w, h float64) (area, perimeter float64) {
    area = w * h
    perimeter = 2 * (w + h)
    return // otomatis return area & perimeter
}
```

## Variadic Function

Menerima **jumlah argumen tak terbatas** menggunakan `...`

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3, 4, 5) // 15
```

`nums` di dalam function bertipe `[]int` (slice).

## Anonymous Function

Function tanpa nama, biasa disimpan di variable.

```go
double := func(x int) int { return x * 2 }
double(7) // 14
```

## Catatan Penting

- Kalau parameter punya tipe yang sama, bisa disingkat: `func add(a, b int)` daripada `func add(a int, b int)`
- Multiple return adalah **pattern paling penting** di Go, terutama `(result, error)`
- Function di Go adalah **first-class citizen** — bisa disimpan di variable, di-pass sebagai argumen, dll
