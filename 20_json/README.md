# 20. JSON Encoding/Decoding

Go punya built-in JSON support via package `encoding/json`.

## Marshal (Struct → JSON)

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{ID: 1, Name: "Ayash", Email: "ayash@mail.com"}
jsonBytes, err := json.Marshal(user)
// {"id":1,"name":"Ayash","email":"ayash@mail.com"}
```

## Unmarshal (JSON → Struct)

```go
jsonStr := `{"id":1,"name":"Ayash"}`
var user User
err := json.Unmarshal([]byte(jsonStr), &user)
```

## Struct Tags

| Tag | Fungsi |
|-----|--------|
| `` `json:"name"` `` | Nama field di JSON |
| `` `json:"name,omitempty"` `` | Skip jika zero value |
| `` `json:"-"` `` | Selalu skip (password, secrets) |

## Di HTTP Handler

```go
// Decode request
json.NewDecoder(r.Body).Decode(&user)

// Encode response
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(user)
```

## Map vs Struct

- **Struct**: Ketika tahu struktur JSON → type-safe, performance bagus
- **Map**: Ketika struktur JSON dinamis → `map[string]interface{}`
