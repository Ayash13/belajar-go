# 5. Looping

Go **hanya punya `for`** — tidak ada `while` atau `do-while`.

## Classic For

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

## While-Style

Hilangkan init dan post statement, jadilah "while":

```go
count := 0
for count < 3 {
    fmt.Println(count)
    count++
}
```

## Infinite Loop

```go
for {
    // jalan terus sampai break
    break
}
```

## Range

Iterasi over **slice**, **map**, **string**, atau **channel**.

### Range over Slice
```go
fruits := []string{"apple", "banana", "mango"}
for i, fruit := range fruits {
    fmt.Printf("[%d] %s\n", i, fruit)
}
```
- `i` = index
- `fruit` = value

### Range over Map
```go
ages := map[string]int{"Ayash": 25, "Budi": 30}
for name, age := range ages {
    fmt.Printf("%s is %d\n", name, age)
}
```

> ⚠️ Urutan iterasi map di Go **tidak dijamin** — bisa berbeda setiap run.

### Skip Index atau Value

```go
for _, fruit := range fruits { } // skip index
for i := range fruits { }        // skip value
```

## Break & Continue

| Keyword | Fungsi |
|---------|--------|
| `break` | Keluar dari loop |
| `continue` | Skip ke iterasi berikutnya |

```go
for i := 0; i < 10; i++ {
    if i == 7 {
        break      // stop di 7
    }
    if i%2 == 0 {
        continue   // skip angka genap
    }
    fmt.Println(i) // cetak: 1, 3, 5
}
```

## Catatan Penting

- Go **hanya punya `for`**, tapi bisa dipakai untuk semua jenis loop
- `range` adalah cara idiomatic Go untuk iterasi collection
- `_` (blank identifier) dipakai untuk mengabaikan value yang tidak dibutuhkan
