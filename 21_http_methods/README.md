# 21. HTTP Methods (GET, POST, PUT, DELETE)

## REST Methods

| Method | Aksi | Contoh Route | Idempotent |
|--------|------|-------------|------------|
| `GET` | Read | `GET /users` | ✅ |
| `POST` | Create | `POST /users` | ❌ |
| `PUT` | Update (full) | `PUT /users/1` | ✅ |
| `PATCH` | Update (partial) | `PATCH /users/1` | ✅ |
| `DELETE` | Delete | `DELETE /users/1` | ✅ |

## Go 1.22+ Routing

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users", listUsers)
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("PUT /users/{id}", updateUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)
```

## Path & Query Parameters

```go
// Path param: /users/123
id := r.PathValue("id") // "123"

// Query param: /users?page=2&limit=10
page := r.URL.Query().Get("page")   // "2"
limit := r.URL.Query().Get("limit") // "10"
```

## RESTful Design

```
GET    /users      → List all
GET    /users/{id} → Get one
POST   /users      → Create
PUT    /users/{id} → Update
DELETE /users/{id} → Delete
```
