# 13. Exported vs Unexported Identifiers

Go tidak punya keyword `public` / `private`. Semua diatur oleh **huruf pertama** dari nama identifier.

## Aturan

| Huruf Pertama | Visibility | Contoh |
|---------------|-----------|--------|
| **Kapital** (A-Z) | Exported — bisa diakses dari package lain | `User`, `NewUser()`, `Name` |
| **Kecil** (a-z) | Unexported — hanya bisa diakses dalam package yang sama | `age`, `getAge()`, `internalConfig` |

## Berlaku Untuk

- **Types**: `User` vs `user`
- **Functions**: `NewUser()` vs `newUser()`
- **Struct Fields**: `Name` vs `name`
- **Methods**: `Greet()` vs `greet()`
- **Variables**: `MaxRetry` vs `maxRetry`
- **Constants**: `Pi` vs `pi`

## Pattern: Getter untuk Unexported Fields

```go
type User struct {
    Name  string  // exported
    age   int     // unexported
}

// Expose via exported getter
func (u User) GetAge() int {
    return u.age
}
```

## Kenapa Penting?

- **Enkapsulasi**: Sembunyikan detail implementasi
- **API Design**: Hanya expose yang perlu dipakai user
- **Safety**: Cegah mutasi state internal secara langsung
