# 7. Structs & Object Modeling

## Deklarasi Struct
Struct digunakan untuk menggabungkan berbagai tipe data menjadi satu entitas.

```go
type User struct {
    ID        int
    FirstName string
    Email     string
    IsActive  bool
}
```

## Inisialisasi
```go
u1 := User{
    ID: 1,
    FirstName: "Ayash",
    Email: "ayash@example.com",
    IsActive: true,
}

// field bisa diakses dengan dot (.)
fmt.Println(u1.FirstName) 
```

## Zero Value
Jika membuat struct tanpa diinisialisasi, field akan mendapat tipe zero value:
```go
var u2 User
// u2.ID = 0, u2.FirstName = "", u2.IsActive = false
```

## Nested Struct (Struct di dalam struct)
```go
type Order struct {
    OrderID string
    User    User
    Total   float64
}
```
Maka aksesnya: `order.User.Email`
