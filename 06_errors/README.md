# 6. Basic Error Handling

## Filosofi Error di Go

Go **tidak punya try/catch/exception**. Sebagai gantinya, error di-return sebagai **value biasa**.

```go
result, err := doSomething()
if err != nil {
    // handle error
}
```

Ini adalah **pattern paling fundamental** di Go.

## Membuat Error

### `errors.New()`
```go
import "errors"

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

- `nil` artinya **tidak ada error** (sukses)
- Return `error` selalu jadi **parameter terakhir**

### `fmt.Errorf()` — Error dengan format
```go
return fmt.Errorf("invalid age: %d", age)
```

## Check Error

```go
result, err := divide(10, 0)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result)
```

> **Selalu check error.** Mengabaikan error adalah anti-pattern di Go.

## Custom Error Type

Implement interface `error` (cukup punya method `Error() string`):

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

Penggunaan:
```go
func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{Field: "age", Message: "cannot be negative"}
    }
    return nil
}
```

## Kenapa Go Pakai Ini?

| Bahasa Lain (try/catch) | Go (error as value) |
|--------------------------|---------------------|
| Error bisa "tersembunyi" | Error **eksplisit** — harus di-handle |
| Bisa lupa catch | Compiler peringatkan kalau error diabaikan |
| Exception mahal secara performa | Error adalah value biasa, ringan |

## Catatan Penting

- `error` adalah **interface** bawaan Go: `type error interface { Error() string }`
- Pattern `if err != nil` akan kamu tulis **sangat sering** — itu normal di Go
- `nil` pada tipe error artinya **sukses / tidak ada error**
- Selalu return error sebagai **parameter terakhir**: `func foo() (result, error)`
