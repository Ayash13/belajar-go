# 10. Interfaces as Contracts

Interface adalah sekumpulan definisi method. Type apapun yang memiliki method-method tersebut secara implisit meng-implementasikan interface-nya (tanpa keyword `implements` seperti di Java).

## Deklarasi
```go
type Shape interface {
    Area() float64
}
```

## Implementasi
Type `Square` akan otomatis masuk kategori `Shape` asalkan ia memiliki method `Area() float64`.

```go
type Square struct {
    Side float64
}

func (s Square) Area() float64 {
    return s.Side * s.Side
}
```

## Kegunaan Utama
Membuat kode menjadi lebih fleksibel/polymorphic:
```go
func PrintArea(s Shape) {
    fmt.Println(s.Area())
}
```
Fungsi di atas bisa menerima tipe data apapun: Square, Circle, Triangle, dsb selama mengimplementasikan method `Area()`.
