# Challenge 3 — REST API Bank (Separation of Concerns)

REST API untuk manajemen rekening bank dan transaksi transfer, dibangun dengan prinsip **Separation of Concerns (SOC)**.

## 📁 Struktur Project

```
challenge_3/
├── main.go                  # Entry point, dependency injection
├── database/
│   └── db.go                # Koneksi PostgreSQL + auto-migrate tabel
├── entity/
│   ├── account.go           # Struct Account
│   └── transaction.go       # Struct Transaction
├── dto/
│   ├── base_response.go     # Struct response standar
│   ├── account_dto.go       # Request/Response untuk Account
│   └── transaction_dto.go   # Request/Response untuk Transaction
├── repository/
│   ├── account_repository.go     # CRUD database Account
│   └── transaction_repository.go # CRUD database Transaction
├── service/
│   └── account_service.go   # Business logic
├── handler/
│   ├── account_handler.go   # HTTP handler (controller)
│   └── account_route.go     # Route mapping
├── server/
│   ├── helper.go            # Helper untuk route pattern
│   └── http.go              # Middleware (Content-Type, 404)
└── .env                     # Environment variables
```

## 🏗️ Arsitektur

```
Request → Handler → Service → Repository → Database
                        ↕
                       DTO ↔ Entity
```

| Layer          | Tanggung Jawab                                     |
|----------------|-----------------------------------------------------|
| **Handler**    | Menerima HTTP request, validasi input, kirim response |
| **Service**    | Business logic (validasi saldo, transfer, dll)       |
| **Repository** | Akses database (query SQL)                           |
| **Entity**     | Representasi tabel database                          |
| **DTO**        | Struktur data request/response API                   |

## ⚙️ Setup

### 1. Buat Database PostgreSQL

```sql
CREATE DATABASE bank;
```

### 2. Konfigurasi Environment

Buat file `.env` di folder `challenge_3/`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=bank
DB_SSLMODE=disable
```

> Tabel `accounts` dan `transactions` akan dibuat otomatis saat aplikasi pertama kali dijalankan.

### 3. Jalankan Aplikasi

```bash
cd challenge_3
go run .
```

Server berjalan di `http://localhost:8080`

## 📋 Database Schema

### accounts

| Column         | Type           | Constraint                    |
|----------------|----------------|-------------------------------|
| id             | UUID           | PRIMARY KEY, auto-generated   |
| account_holder | VARCHAR(255)   | NOT NULL                      |
| balance        | NUMERIC(15,2)  | NOT NULL, DEFAULT 0           |
| created_at     | TIMESTAMP      | NOT NULL, DEFAULT NOW()       |
| updated_at     | TIMESTAMP      | NOT NULL, DEFAULT NOW()       |

### transactions

| Column          | Type           | Constraint                           |
|-----------------|----------------|--------------------------------------|
| id              | UUID           | PRIMARY KEY, auto-generated          |
| from_account_id | UUID           | NOT NULL, FK → accounts(id)          |
| to_account_id   | UUID           | NOT NULL, FK → accounts(id)          |
| amount          | NUMERIC(15,2)  | NOT NULL                             |
| created_at      | TIMESTAMP      | NOT NULL, DEFAULT NOW()              |

## 🔌 API Endpoints

Base URL: `http://localhost:8080`

### 1. Create Account

```
POST /accounts
```

**Request Body:**
```json
{
    "account_holder": "John Doe",
    "balance": 50000
}
```

**Response (201):**
```json
{
    "code": 201,
    "status": "success",
    "message": "Account created successfully",
    "data": {
        "id": "uuid-here",
        "account_holder": "John Doe",
        "balance": 50000,
        "created_at": "2026-03-25T11:05:43.918645Z",
        "updated_at": "2026-03-25T11:05:43.918645Z"
    }
}
```

---

### 2. Get All Accounts

```
GET /accounts
```

**Response (200):**
```json
{
    "code": 200,
    "status": "success",
    "data": [
        {
            "id": "uuid-1",
            "account_holder": "John Doe",
            "balance": 50000,
            "created_at": "...",
            "updated_at": "..."
        }
    ]
}
```

---

### 3. Get Account By ID

```
GET /accounts/{id}
```

**Response (200):**
```json
{
    "code": 200,
    "status": "success",
    "data": {
        "id": "uuid-here",
        "account_holder": "John Doe",
        "balance": 50000,
        "created_at": "...",
        "updated_at": "..."
    }
}
```

**Response (404):**
```json
{
    "code": 404,
    "status": "error",
    "message": "account not found"
}
```

---

### 4. Update Account

```
PUT /accounts/{id}
```

**Request Body:**
```json
{
    "account_holder": "John Updated",
    "balance": 75000
}
```

**Response (200):**
```json
{
    "code": 200,
    "status": "success",
    "message": "Account updated successfully",
    "data": {
        "id": "uuid-here",
        "account_holder": "John Updated",
        "balance": 75000,
        "created_at": "...",
        "updated_at": "..."
    }
}
```

---

### 5. Delete Account

```
DELETE /accounts/{id}
```

**Response (200):**
```json
{
    "code": 200,
    "status": "success",
    "message": "Account deleted successfully"
}
```

> ⚠️ Tidak bisa menghapus akun yang masih punya riwayat transaksi (foreign key constraint).

---

### 6. Transfer

```
POST /transfer
```

**Request Body:**
```json
{
    "from_account_id": "uuid-sender",
    "to_account_id": "uuid-receiver",
    "amount": 10000
}
```

**Response (200):**
```json
{
    "code": 200,
    "status": "success",
    "data": {
        "message": "Transfer successful",
        "transaction": {
            "id": "uuid-transaction",
            "from_account_id": "uuid-sender",
            "to_account_id": "uuid-receiver",
            "amount": 10000,
            "created_at": "..."
        },
        "from_account": {
            "id": "uuid-sender",
            "account_holder": "John Doe",
            "balance": 40000,
            "created_at": "...",
            "updated_at": "..."
        },
        "to_account": {
            "id": "uuid-receiver",
            "account_holder": "Jane Doe",
            "balance": 40000,
            "created_at": "...",
            "updated_at": "..."
        }
    }
}
```

**Error — Saldo tidak cukup (400):**
```json
{
    "code": 400,
    "status": "error",
    "message": "insufficient balance for transfer"
}
```

---

### 7. Get Transactions By Account ID

```
GET /accounts/{id}/transactions
```

**Response (200):**
```json
{
    "code": 200,
    "status": "success",
    "data": [
        {
            "id": "uuid-transaction",
            "from_account_id": "uuid-sender",
            "to_account_id": "uuid-receiver",
            "amount": 10000,
            "created_at": "..."
        }
    ]
}
```

## 🛠️ Tech Stack

- **Go** (net/http + Go 1.22 enhanced routing)
- **PostgreSQL** — Database
- **sqlx** — SQL query builder
- **godotenv** — Environment variable loader
- **lib/pq** — PostgreSQL driver

## 📌 Fitur

- ✅ CRUD Rekening Bank
- ✅ Transfer antar rekening (dengan database transaction)
- ✅ Riwayat transaksi per akun
- ✅ Auto-migrate tabel saat startup
- ✅ Middleware Content-Type JSON
- ✅ Custom 404 handler
- ✅ Validasi input request
- ✅ Error handling yang konsisten
