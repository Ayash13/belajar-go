# 23. Middleware — Real Examples

Modul ini menjalankan **HTTP server asli** dan membuat **request nyata** untuk mendemonstrasikan setiap middleware.

## 5 Middleware yang Dicontohkan

### 1. Logging Middleware
Mencatat setiap request: method, path, status code, dan durasi.

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        recorder := &statusRecorder{ResponseWriter: w, status: 200}
        next.ServeHTTP(recorder, r)
        fmt.Printf("%s %s → %d (%v)\n", r.Method, r.URL.Path, recorder.status, time.Since(start))
    })
}
```

### 2. Auth Middleware
Cek `Authorization: Bearer <token>` header. Reject jika missing/invalid.

```go
// No token    → 401 Unauthorized
// Wrong token → 403 Forbidden
// Valid token → next handler
```

### 3. CORS Middleware
Set header `Access-Control-Allow-*` dan handle preflight `OPTIONS` request.

### 4. Rate Limiter
Batasi jumlah request per IP per window waktu.

```go
limiter := newRateLimiter(3, 10*time.Second)
// Request 1-3 → 200 OK
// Request 4+  → 429 Too Many Requests
```

### 5. Recovery Middleware
Tangkap panic di handler, return 500 instead of crash.

```go
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                http.Error(w, "internal error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

## Chaining

```go
// Public:    CORS → Logging → Recovery → Handler
// Protected: CORS → Logging → Recovery → Auth → Handler
// Limited:   Logging → RateLimiter → Handler

handler := corsMiddleware(loggingMiddleware(recoveryMiddleware(authMiddleware(mux))))
```

## Execution Flow

```
Request  → CORS → Logging → Recovery → Auth → Handler
Response ← CORS ← Logging ← Recovery ← Auth ← Handler
```
