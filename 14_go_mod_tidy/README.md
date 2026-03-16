# 14. go mod tidy

`go mod tidy` adalah perintah penting untuk mengelola dependency di Go.

## File Penting

### go.mod
```
module belajar-go

go 1.26.1

require (
    github.com/gorilla/mux v1.8.1
    github.com/lib/pq v1.10.9
)
```

### go.sum
Berisi hash checksum dari setiap dependency. **Jangan edit manual.**

## Perintah

| Perintah | Fungsi |
|----------|--------|
| `go mod init <name>` | Inisialisasi module baru |
| `go mod tidy` | Sinkronkan dependencies (tambah/hapus) |
| `go get <pkg>` | Tambah/update dependency |
| `go mod download` | Download semua dependency |
| `go mod verify` | Verifikasi integrity dependency |

## Apa yang `go mod tidy` Lakukan?

1. **Scan** semua file `.go` untuk mencari import
2. **Tambah** dependency yang dipakai tapi belum di `go.mod`
3. **Hapus** dependency di `go.mod` yang sudah tidak dipakai
4. **Update** `go.sum` dengan checksum yang benar

## Kapan Jalankan?

- Setelah menambah/menghapus import
- Sebelum commit ke git
- Saat clone project baru
- Saat ada error "missing module"
