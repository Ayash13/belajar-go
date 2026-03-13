# 9. Pointer Concepts

Pointer adalah variabel yang menyimpan "alamat memori" dari variabel lain.

## Syntax Dasar
- `&` (address of): untuk mendapatkan alamat memori suatu variable.
- `*` (dereference): untuk membaca/mengubah nilai dari sebuah alamat memori.

```go
x := 10
p := &x          // p sekarang adalah pointer yang menunjuk ke lokasi memori x
fmt.Println(p)   // misal: 0xc0000180b0
fmt.Println(*p)  // 10

*p = 20          // x ikut berubah menjadi 20
```

## Passing Pointer ke Fungsi
Secara default parameter di Go adalah **pass-by-value** (di-copy). Jika ingin fungsi bisa mengubah nilai asli, pass menggunakan pointer.

```go
func passByPointer(ptr *int) {
    *ptr = 99
}

num := 5
passByPointer(&num)
// num sekarang 99
```
