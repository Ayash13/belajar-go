# 18. Basic HTTP Server

Go punya built-in HTTP server di package `net/http`. Tidak perlu framework eksternal.

## Minimal Server

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })

    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```

## Komponen Utama

| Komponen | Fungsi |
|----------|--------|
| `http.HandleFunc` | Daftarkan handler untuk route tertentu |
| `http.ResponseWriter` | Interface untuk menulis response |
| `*http.Request` | Struct berisi data request (method, URL, body, headers) |
| `http.ListenAndServe` | Start server pada address tertentu |

## Custom Server (Recommended)

```go
srv := &http.Server{
    Addr:         ":8080",
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
srv.ListenAndServe()
```
