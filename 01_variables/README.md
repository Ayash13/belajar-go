# 1. Variables & Data Types

## Deklarasi Variable

Go punya **3 cara** deklarasi variable:

### `var` keyword
```go
var name string = "Ayash"
var age int = 25
```
Dipakai kalau mau **eksplisit** tentukan tipe data.

### Short declaration `:=`
```go
score := 95.5
city := "Jakarta"
```
Cara paling **sering dipakai**. Go otomatis menebak tipe datanya (type inference). Hanya bisa dipakai **di dalam function**.

### Multiple declaration
```go
var a, b, c int = 1, 2, 3
```

## Tipe Data Dasar

| Tipe | Contoh | Keterangan |
|------|--------|------------|
| `string` | `"hello"` | Teks |
| `int` | `42` | Bilangan bulat |
| `float64` | `3.14` | Bilangan desimal |
| `bool` | `true` / `false` | Boolean |

## Zero Values

Kalau variable dideklarasi tapi **tidak diberi nilai**, Go otomatis kasih **zero value**:

| Tipe | Zero Value |
|------|-----------|
| `int` | `0` |
| `float64` | `0.0` |
| `string` | `""` (string kosong) |
| `bool` | `false` |

```go
var defaultInt int       // 0
var defaultString string // ""
var defaultBool bool     // false
```

## Cek Tipe Data

Gunakan `%T` di `fmt.Printf`:

```go
fmt.Printf("tipe: %T\n", score) // tipe: float64
```

## Catatan Penting

- `:=` hanya bisa di dalam function, tidak bisa di level package
- Variable yang dideklarasi tapi **tidak dipakai** akan error di Go
- Go adalah **statically typed** — tipe data tidak bisa berubah setelah dideklarasi
