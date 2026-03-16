# 12. Package System

Di Go, kode diorganisasi dalam **packages**. Setiap folder adalah satu package.

## Aturan Package

1. Semua file `.go` dalam satu folder harus punya `package` name yang sama
2. Import menggunakan path relatif dari module name di `go.mod`
3. Package `main` adalah entry point — hanya `main` yang bisa punya `func main()`

## Struktur

```
12_packages/
├── packages.go          # package "packages" (entry point modul ini)
├── mathutil/
│   └── mathutil.go      # package "mathutil"
└── stringutil/
    └── stringutil.go    # package "stringutil"
```

## Cara Import

```go
import (
    "belajar-go/12_packages/mathutil"
    "belajar-go/12_packages/stringutil"
)

result := mathutil.Add(1, 2)
reversed := stringutil.Reverse("Go")
```

## Tips

- Nama package sebaiknya **pendek dan deskriptif** (`mathutil`, bukan `my_math_utility_helpers`)
- Hindari nama package yang terlalu umum seperti `util` atau `common`
- Package name tidak harus sama dengan folder name, tapi **konvensi** di Go membuatnya sama
