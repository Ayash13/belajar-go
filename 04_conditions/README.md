# 4. If/Else & Switch

## If/Else

```go
if score >= 90 {
    fmt.Println("A")
} else if score >= 80 {
    fmt.Println("B")
} else {
    fmt.Println("C")
}
```

> **Catatan**: Di Go, **tidak pakai kurung `()`** di kondisi. Tapi **kurung kurawal `{}` wajib**.

## If dengan Short Statement

Go bisa deklarasi variable **di dalam if**. Variable tersebut hanya bisa diakses di scope if/else.

```go
if x := 10; x > 5 {
    fmt.Println(x, "lebih besar dari 5")
}
// x tidak bisa diakses di sini
```

Ini berguna untuk:
```go
if err := doSomething(); err != nil {
    // handle error
}
```

## Switch

Lebih bersih dari banyak `if/else`.

```go
switch day {
case "Monday", "Tuesday", "Wednesday":
    fmt.Println("weekday")
case "Saturday", "Sunday":
    fmt.Println("weekend")
default:
    fmt.Println("unknown")
}
```

### Perbedaan dengan bahasa lain:
- **Tidak perlu `break`** — Go otomatis break di setiap case
- Bisa punya **multiple values** di satu case: `case "Monday", "Tuesday":`
- Kalau mau lanjut ke case berikutnya, pakai `fallthrough`

## Switch Tanpa Expression

Bisa dipakai seperti if/else chain:

```go
switch {
case temp > 30:
    fmt.Println("hot!")
case temp > 20:
    fmt.Println("warm")
default:
    fmt.Println("cold")
}
```

## Catatan Penting

- `{` **harus** di baris yang sama dengan `if`/`else`/`switch` — Go tidak mengizinkan di baris baru
- Switch di Go **tidak fallthrough** secara default (beda dengan C/Java)
- Gunakan switch kalau ada **3+ kondisi** — lebih readable daripada if/else chain
