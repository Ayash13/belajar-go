# Practice 2: PostgreSQL CRUD API

REST API sederhana untuk manage tasks, terhubung ke PostgreSQL.

## Prerequisites

1. PostgreSQL harus sudah terinstall dan running
2. Buat database:
```sql
CREATE DATABASE belajar_go;
```

## Setup & Run

```bash
# Install dependency
cd practice_02_postgres_crud
go mod tidy

# Jalankan (default: localhost:5432, user=postgres, password=postgres)
go run main.go

# Atau dengan custom connection string
DATABASE_URL="host=localhost port=5432 user=postgres password=yourpassword dbname=belajar_go sslmode=disable" go run main.go
```

Tabel `tasks` akan otomatis dibuat saat server start.

## API Endpoints

| Method | Endpoint | Body | Keterangan |
|--------|----------|------|------------|
| `GET` | `/tasks` | — | List semua tasks |
| `GET` | `/tasks/1` | — | Get task by ID |
| `POST` | `/tasks` | `{"title":"...","description":"..."}` | Create task |
| `PUT` | `/tasks/1` | `{"title":"...","completed":true}` | Update task |
| `DELETE` | `/tasks/1` | — | Delete task |

## Contoh Request (curl)

```bash
# Create
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Belajar Go","description":"Selesaikan semua modul"}'

# List
curl http://localhost:8080/tasks

# Get by ID
curl http://localhost:8080/tasks/1

# Update
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'

# Delete
curl -X DELETE http://localhost:8080/tasks/1
```

## Response Format

```json
{
  "success": true,
  "message": "task created",
  "data": {
    "id": 1,
    "title": "Belajar Go",
    "description": "Selesaikan semua modul",
    "completed": false,
    "created_at": "2026-03-16T13:00:00Z",
    "updated_at": "2026-03-16T13:00:00Z"
  }
}
```

## Konsep yang Dipakai

| Modul | Konsep |
|-------|--------|
| 12 | Package system (imports) |
| 13 | Exported types & functions |
| 14 | go mod tidy (dependency management) |
| 15 | Separation of concerns (model/repo/handler) |
| 16 | Error wrapping |
| 17 | Database integration (sql.Open, Query, Scan) |
| 18 | HTTP server (ListenAndServe) |
| 19 | Handlers (HandleFunc) |
| 20 | JSON encoding/decoding |
| 21 | HTTP methods (GET/POST/PUT/DELETE) |
| 22 | Status codes (200, 201, 400, 404, 500) |
| 23 | Middleware (logging, CORS) |

## Database Schema

```sql
CREATE TABLE tasks (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    completed   BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);
```
