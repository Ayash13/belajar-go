# Practice 1: Simple API Fetch

Mari mempraktekkan materi sebelumnya (Struct, Error Handling) dengan membuat HTTP request sederhana.

## Yang Dipelajari
1. Melakukan HTTP `GET` request menggunakan package `net/http`
2. Membaca response body menggunakan package `io`
3. Parsing JSON (Unmarshal) menjadi Struct Go menggunakan package `encoding/json`
4. Menggunakan struct tag `json:"field_name"`
5. `defer` statement untuk memastikan koneksi / file ditutup

## Struct Tags
Karena JSON dari API menggunakan camelCase (misal `userId`), sedangkan field di Go umumnya menggunakan PascalCase (misal `UserID`), kita butuh struct tags agar Go tahu cara memetakannya.

```go
type Todo struct {
    UserID int `json:"userId"`
}
```

## Defer
Keyword `defer` menjadwalkan pemanggilan fungsi untuk dieksekusi tepat sebelum fungsi yang melingkupinya return. Sangat berguna untuk `Close()` resource.

```go
resp, _ := http.Get("URL")
defer resp.Body.Close() // Pasti dieksekusi di akhir fungsi
```
