# 19. Handlers

Handler adalah fungsi yang memproses HTTP request dan mengirim response.

## 3 Cara Membuat Handler

### 1. Inline Function
```go
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello!")
})
```

### 2. Named Function
```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello!")
}
http.HandleFunc("/hello", helloHandler)
```

### 3. Struct (implement `http.Handler`)
```go
type GreetHandler struct {
    Greeting string
}

func (h *GreetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, h.Greeting)
}

http.Handle("/greet", &GreetHandler{Greeting: "Halo!"})
```

## Request Data

| Property | Contoh | Keterangan |
|----------|--------|------------|
| `r.Method` | `GET` | HTTP method |
| `r.URL.Path` | `/users/123` | URL path |
| `r.URL.Query()` | `?name=Go` | Query params |
| `r.Header.Get("Authorization")` | `Bearer xxx` | Request headers |
| `r.Body` | `io.ReadCloser` | Request body |
| `r.PathValue("id")` | `123` | Path params (Go 1.22+) |

## Response

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK) // 200
w.Write([]byte(`{"message":"ok"}`))
```
