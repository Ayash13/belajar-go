# 2. Constants, Iota & Operators

## Constants

Nilai yang **tidak bisa diubah** setelah dideklarasi.

```go
const Pi = 3.14
const AppName = "BelajarGo"
```

Kalau coba reassign `Pi = 3.15` → **compile error**.

## Iota

`iota` adalah **auto-increment counter** yang dimulai dari `0`. Biasa dipakai untuk enum.

```go
const (
    Sunday    = iota // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
)
```

Setiap baris baru dalam block `const`, `iota` otomatis naik 1.

### Trik Iota

```go
// Skip value
const (
    _ = iota // 0 (skip)
    KB = 1 << (10 * iota) // 1024
    MB                     // 1048576
    GB                     // 1073741824
)
```

## Operators

### Arithmetic (Aritmatika)
| Operator | Contoh | Hasil |
|----------|--------|-------|
| `+` | `10 + 3` | `13` |
| `-` | `10 - 3` | `7` |
| `*` | `10 * 3` | `30` |
| `/` | `10 / 3` | `3` (integer division) |
| `%` | `10 % 3` | `1` (sisa bagi) |

> **Catatan**: `10 / 3` hasilnya `3` bukan `3.33` karena keduanya `int`. Kalau mau desimal, pakai `float64`.

### Comparison (Perbandingan)
| Operator | Arti |
|----------|------|
| `==` | sama dengan |
| `!=` | tidak sama |
| `>` | lebih besar |
| `<` | lebih kecil |
| `>=` | lebih besar atau sama |
| `<=` | lebih kecil atau sama |

### Logical (Logika)
| Operator | Arti | Contoh |
|----------|------|--------|
| `&&` | AND | `true && false` → `false` |
| `\|\|` | OR | `true \|\| false` → `true` |
| `!` | NOT | `!true` → `false` |
