# 16. Error Wrapping

Go 1.13+ memperkenalkan error wrapping — menambahkan konteks ke error sambil menjaga error chain.

## `fmt.Errorf` dengan `%w`

```go
// Wrap error dengan konteks
return fmt.Errorf("repo.GetUser(id=%d): %w", id, err)
// Output: repo.GetUser(id=99): not found

// ⚠️ %v memutus chain, %w mempertahankan chain
```

## `errors.Is` — Cek Error di Chain

```go
var ErrNotFound = errors.New("not found")

err := fmt.Errorf("service: %w", ErrNotFound)
errors.Is(err, ErrNotFound) // true ✅
```

## `errors.As` — Extract Error Type

```go
type AppError struct {
    Code    int
    Message string
    Err     error
}

var target *AppError
if errors.As(err, &target) {
    fmt.Println(target.Code)    // 404
    fmt.Println(target.Message) // "user not found"
}
```

## Custom Error dengan `Unwrap()`

```go
func (e *AppError) Error() string {
    return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
    return e.Err  // memungkinkan errors.Is/As menelusuri chain
}
```

## `%w` vs `%v`

| | `%w` | `%v` |
|---|------|------|
| Error chain | ✅ Utuh | ❌ Putus |
| `errors.Is` | ✅ Bisa match | ❌ Tidak match |
| `errors.As` | ✅ Bisa extract | ❌ Tidak bisa |
| Kapan pakai | Default | Saat ingin hide internal error |
