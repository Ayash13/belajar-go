# 17. Database Integration (PostgreSQL)

Modul ini menggunakan **koneksi PostgreSQL asli** dengan dua library populer: **sqlx** dan **GORM**.

## Prerequisites

```bash
# Pastikan PostgreSQL running, lalu buat database:
psql -U postgres -c "CREATE DATABASE belajar_go;"
```

## Library

### sqlx — Enhanced `database/sql`
```bash
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
```

Sqlx memperluas `database/sql` standar dengan fitur tambahan:
- `db.Get()` → scan single row langsung ke struct
- `db.Select()` → scan multiple rows langsung ke slice
- `db.NamedExec()` → query dengan named parameters
- `sqlx.In()` → helper untuk `WHERE IN` query

### GORM — Full-featured ORM
```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

GORM adalah ORM lengkap:
- `db.Create()` → INSERT
- `db.First()` / `db.Find()` → SELECT
- `db.Update()` / `db.Updates()` → UPDATE
- `db.Delete()` → DELETE
- `db.AutoMigrate()` → auto create/alter table

## Struct Tags

```go
// sqlx — menggunakan `db` tag
type Task struct {
    ID    int    `db:"id"`
    Title string `db:"title"`
}

// GORM — menggunakan `gorm` tag
type Task struct {
    ID    uint   `gorm:"primaryKey"`
    Title string `gorm:"size:255;not null"`
}
```

## Perbandingan

| | sqlx | GORM |
|---|------|------|
| Style | SQL-first | ORM (method chaining) |
| Query | Raw SQL | Chain methods |
| Migration | Manual | AutoMigrate |
| Learning curve | Mudah | Sedang |
| Control | Tinggi | Sedang |
| Performance | Lebih cepat | Sedikit overhead |
| Cocok untuk | Simple projects | Complex projects |

## Connection String

```
host=localhost port=5432 user=postgres password=postgres dbname=belajar_go sslmode=disable
```

Atau set environment variable:
```bash
export DATABASE_URL="host=localhost port=5432 user=postgres password=yourpassword dbname=belajar_go sslmode=disable"
```
