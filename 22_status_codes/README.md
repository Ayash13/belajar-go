# 22. HTTP Status Codes

Setiap HTTP response punya status code yang menunjukkan hasil request.

## Cheat Sheet

| Code | Name | Kapan Pakai |
|------|------|------------|
| **200** | OK | Request berhasil |
| **201** | Created | Resource baru dibuat (POST) |
| **204** | No Content | Berhasil tapi tanpa body (DELETE) |
| **400** | Bad Request | Request tidak valid |
| **401** | Unauthorized | Belum login |
| **403** | Forbidden | Tidak punya akses |
| **404** | Not Found | Resource tidak ada |
| **405** | Method Not Allowed | HTTP method salah |
| **409** | Conflict | Data konflik (duplicate) |
| **422** | Unprocessable Entity | Validasi gagal |
| **500** | Internal Server Error | Bug di server |

## Di Go

```go
// Selalu gunakan constant
w.WriteHeader(http.StatusOK)           // 200
w.WriteHeader(http.StatusCreated)      // 201
w.WriteHeader(http.StatusNotFound)     // 404

// Helper untuk error response
http.Error(w, "not found", http.StatusNotFound)
```

## ⚠️ Aturan Penting

1. `WriteHeader()` harus dipanggil **sebelum** `Write()`
2. `WriteHeader()` hanya bisa dipanggil **sekali**
3. Gunakan **constant** (`http.StatusOK`), bukan angka (`200`)
