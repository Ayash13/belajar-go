# 11. Dependency Injection

Dependency Injection (DI) adalah menyuntikkan ketergantungan (dependencies) ke sebuah struct/fungsi dari luar, bukan menginstansiasi di dalam, membuat sistem gampang di-test.

## Kenapa menggunakan DI?
Di Go, DI biasanya menggunakan **Interface**.

```go
// Interface (Contract)
type Logger interface {
    Log(message string)
}

// Service yang bergantung pada Logger
type UserService struct {
    logger Logger
}

// Constructor Injection
func NewUserService(l Logger) *UserService {
    return &UserService{logger: l}
}
```

Dengan design ini, tipe `Logger` yang masuk bisa berupa:
- `ConsoleLogger`: Cetak di terminal
- `FileLogger`: Tulis ke file teks
- `MockLogger`: Untuk unit testing (tanpa perlu benar-benar print/tulis file)
