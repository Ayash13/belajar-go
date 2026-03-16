# 15. Separation of Concerns

Separation of Concerns (SoC) memisahkan kode berdasarkan tanggung jawab. Setiap bagian hanya fokus pada satu hal.

## Layer Architecture

```
┌─────────────────┐
│   Presenter      │  ← Tampilan / Output
├─────────────────┤
│   Service        │  ← Business Logic
├─────────────────┤
│   Repository     │  ← Data Access (CRUD)
├─────────────────┤
│   Model          │  ← Data Structure
└─────────────────┘
```

## Contoh

```go
// MODEL — hanya data
type Product struct {
    ID    int
    Name  string
    Price float64
}

// REPOSITORY — akses data
type ProductRepository interface {
    FindAll() []Product
    FindByID(id int) (Product, bool)
}

// SERVICE — business logic
type ProductService struct {
    repo ProductRepository
}

func (s *ProductService) IsExpensive(p Product) bool {
    return p.Price > 1000000
}
```

## Kenapa Penting?

| Tanpa SoC | Dengan SoC |
|-----------|-----------|
| Semua logika campur di satu file | Setiap layer punya file sendiri |
| Sulit di-test | Mudah di-test (mock repository) |
| Ganti database = rewrite semua | Ganti database = buat repo baru |
| Sulit dibaca | Kode terstruktur dan jelas |

## Di Proyek Nyata

```
project/
├── models/       # Struct definitions
├── repository/   # Database queries
├── service/      # Business logic
├── handler/      # HTTP handlers (presenter)
└── main.go       # Wiring / setup
```
