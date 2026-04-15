# Challenge 3 — REST API Bank (Separation of Concerns)

REST API untuk manajemen rekening bank dan transaksi transfer, dibangun dengan prinsip **Separation of Concerns (SOC)**.

## 📁 Struktur Project

```
challenge_3/
├── main.go                     # Entry point, Dependency Injection & Middleware Wiring
├── .env                        # Konfigurasi Environment (DB, Redis, etc.)
├── docker-compose.yml          # Orkestrasi Stack (App, DB, Redis, Grafana Stack)
├── Dockerfile                  # Containerisasi Aplikasi Go
│
├── database/                   # Layer Koneksi Data
│   ├── db.go                   # Koneksi & Auto-Migrate PostgreSQL
│   └── redis.go                # Koneksi Redis Client
│
├── entity/                     # Domain Models (Database Schema)
│   ├── account.go              # Model Account
│   └── transaction.go          # Model Transaction
│
├── dto/                        # Data Transfer Object (API Contracts)
│   ├── base_response.go        # Standard Response Format
│   ├── account_dto.go          # request/Response DTO Account
│   └── transaction_dto.go      # Request/Response DTO Transaction
│
├── repository/                 # Layer Akses Data (SQL Queries)
│   ├── account_repository.go    # Query CRUD Account
│   └── transaction_repository.go # Query Transaksi & Transfer Logic
│
├── service/                    # Layer Bisnis Logic
│   └── account_service.go      # Validasi saldo, transfer validation, logic bisnis
│
├── handler/                    # Layer Transport (HTTP)
│   ├── account_handler.go      # Controller/Handler Logic
│   ├── account_route.go        # Mapping HTTP Endpoints
│   └── account_handler_test.go # Unit Test untuk API Layer
│
├── middleware/                 # Cross-Cutting Concerns
│   ├── core.go                 # Root Middleware (Content-Type, Chain Composer)
│   ├── tracing.go              # OpenTelemetry Tracing instrumentation
│   ├── metrics.go              # Prometheus Metrics recording
│   ├── request_logger.go       # Structured Logging (Zap)
│   ├── rate_limit.go           # Redis-based Rate Limiting (Fixed Window)
│   ├── idempotency.go          # Redis-based Idempotency check
│   └── cache.go                # Redis Caching logic untuk GET requests
│
├── telemetry/                  # Konfigurasi Observability
│   ├── tracer.go               # Setup OpenTelemetry (Tempo Exporter)
│   └── metrics.go              # Setup Prometheus Metrics (Counter, Histogram)
│
├── logger/                     # Logging Configuration
│   └── logger.go               # Zap Logger Initialization (Structured Logging)
│
├── config/                     # Configuration Management
│   ├── prometheus.yaml         # Scrape Job Metrics
│   ├── tempo.yaml              # Backend Traces Config
│   ├── promtail.yaml           # Log Collector Config
│   └── grafana-datasources.yaml # Auto-provisioning Grafana
│
└── server/                     # HTTP Server Utilities
    └── helper.go               # Route & Pattern matching helpers
```

## 🛡️ Middleware Chain Order

Setiap request masuk melalui rantai (chain) middleware dengan urutan berikut:

1.  **OpenTelemetry**: Inisialisasi span dan injector context.
2.  **Prometheus Metrics**: Mulai timer untuk record latensi dan hitung total request.
3.  **Request Logger**: Mencatat log masuk (Method, Path, IP).
4.  **Content-Type JSON**: Memaksa header `application/json`.
5.  **Rate Limiter**: Cek IP klien apakah melebihi batas quota (Redis).
6.  **Idempotency**: Cek header `Idempotency-Key` untuk mencegah transaksi ganda (Redis).
7.  **Timeout**: Membatasi eksekusi request (misal: 10 detik).
8.  **Handler/Cache**: Jika GET, cek cache di Redis. Jika tidak ada, eksekusi service.

## 📊 Observability (Three Pillars)

Sistem ini menerapkan **Full Stack Observability** yang dapat dipantau melalui Grafana (`localhost:3000`):

| Pillar      | Provider        | Tujuan                                               | Dashboard Path       |
| :---------- | :-------------- | :--------------------------------------------------- | :------------------- |
| **Metrics** | Prometheus      | Pantau throughput (RPS), Error Rate, dan Latensi.    | Explore → Prometheus |
| **Logs**    | Loki + Promtail | Analisis log error terstruktur (JSON).               | Explore → Loki       |
| **Traces**  | OpenTelemetry   | Visualisasi alur request antar komponen (SDK level). | Explore → Tempo      |

## 🏗️ Arsitektur

```
Request → Handler → Service → Repository → Database
                        ↕
                       DTO ↔ Entity
```

| Layer          | Tanggung Jawab                                        |
| -------------- | ----------------------------------------------------- |
| **Handler**    | Menerima HTTP request, validasi input, kirim response |
| **Service**    | Business logic (validasi saldo, transfer, dll)        |
| **Repository** | Akses database (query SQL)                            |
| **Entity**     | Representasi tabel database                           |
| **DTO**        | Struktur data request/response API                    |

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

Terdapat dua cara untuk menjalankan aplikasi ini: secara manual atau menggunakan Docker.

#### Menggunakan Docker Compose (Direkomendasikan)

Pastikan Docker telah terinstall di sistem Anda, jalankan perintah berikut:

```bash
cd challenge_3
docker-compose up -d --build
```

> Perintah ini akan mensetup database `bank` di dalam container `postgres` secara otomatis bersamaan dengan menjalankan server HTTP.

#### Menjalankan Secara Manual

```bash
cd challenge_3
go run .
```

Server berjalan di `http://localhost:8080`

## 📋 Database Schema

### accounts

| Column         | Type          | Constraint                  |
| -------------- | ------------- | --------------------------- |
| id             | UUID          | PRIMARY KEY, auto-generated |
| account_holder | VARCHAR(255)  | NOT NULL                    |
| balance        | NUMERIC(15,2) | NOT NULL, DEFAULT 0         |
| created_at     | TIMESTAMP     | NOT NULL, DEFAULT NOW()     |
| updated_at     | TIMESTAMP     | NOT NULL, DEFAULT NOW()     |

### transactions

| Column          | Type          | Constraint                  |
| --------------- | ------------- | --------------------------- |
| id              | UUID          | PRIMARY KEY, auto-generated |
| from_account_id | UUID          | NOT NULL, FK → accounts(id) |
| to_account_id   | UUID          | NOT NULL, FK → accounts(id) |
| amount          | NUMERIC(15,2) | NOT NULL                    |
| created_at      | TIMESTAMP     | NOT NULL, DEFAULT NOW()     |

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

**HTTP Headers:**

```
Content-Type: application/json
Authorization: Bearer <token>
Authorization-Customer: Bearer <customer_token>
X-TIMESTAMP: 2020-12-21T10:30:24+07:00
X-SIGNATURE: 85be817c55...
ORIGIN: www.hostname.com
X-PARTNER-ID: 821508239...
X-EXTERNAL-ID: 418075533...
X-IP-ADDRESS: 172.24.281.24
X-DEVICE-ID: 09864ADCASA
X-LATITUDE: -6.1617169
X-LONGITUDE: 106.6643946
CHANNEL-ID: 95221
```

**Request Body:**

```json
{
   "partnerReferenceNo":"2020102900000000000021",
   "amount":{
      "value":"10000.00",
      "currency":"IDR"
   },
   "beneficiaryAccountNo":"uuid-receiver",
   "beneficiaryEmail":"yories.yolanda@work.bri.co.id",
   "currency":"IDR",
   "customerReference":"10052019",
   "feeType":"BEN",
   "remark":"remark test",
   "sourceAccountNo":"uuid-sender",
   "transactionDate":"2019-07-03T12:08:56+07:00",
   "originatorInfos":[{
      "originatorCustomerNo":"999901000003300",
      "originatorCustomerName":"Hafizh",
      "originatorBankCode":"002"
   }],
   "additionalInfo":{
      "deviceId":"12345679237",
      "channel":"mobilephone"
   }
}
```

**Response (200 Success):**

```json
{
   "responseCode": "2001700",
   "responseMessage": "Successful",
   "referenceNo": "uuid-transaction",
   "partnerReferenceNo": "2020102900000000000021",
   "amount": {
      "value": "10000.00",
      "currency": "IDR"
   },
   "beneficiaryAccountNo": "uuid-receiver",
   "currency": "IDR",
   "customerReference": "10052019",
   "sourceAccount": "uuid-sender",
   "transactionDate": "2019-07-03T12:08:56+07:00",
   "originatorInfos": [ ... ],
   "additionalInfo": { ... }
}
```

**Error — Saldo tidak cukup (403):**

```json
{
  "responseCode": "4031714",
  "responseMessage": "Insufficient Funds"
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

## 📊 Observability (Three Pillars)

Proyek ini menerapkan konsep **Three Pillars of Observability** agar sistem dapat dipantau secara menyeluruh.

### 1. Metrics (Prometheus)

Digunakan untuk menjawab pertanyaan: _"Seberapa sibuk server kita dan berapa lama waktu responnya?"_

- **Path**: `GET /metrics`
- **Metrik Utama (HTTP)**: `bank_api_http_requests_total` dan `bank_api_http_request_duration_seconds`.
- **Metrik Bisnis (Transfer)**: `bank_api_business_transactions_total`, `bank_api_business_transfer_amount_total_usd`, dan `bank_api_business_transfer_amount_distribution`.
- **Metrik Sistem**: `bank_api_business_accounts_created_total`.

### 2. Logs (Loki + Promtail)

Digunakan untuk menjawab pertanyaan: _"Apa yang sebenarnya terjadi saat error muncul?"_

- **Teknologi**: Zap Logger (JSON format) dikumpulkan secara otomatis oleh Promtail dan dikirim ke Loki.
- **Trace ID**: Setiap log menyertakan `trace_id` sehingga Anda bisa mencocokkan log dengan trace tertentu.

### 3. Traces (OpenTelemetry + Tempo)

Digunakan untuk menjawab pertanyaan: _"Dimana bottleneck pemrosesan request ini?"_

- **Teknologi**: OpenTelemetry SDK melakukan instrumen pada level Handler, Service, dan Database.
- **Visualisasi**: Traces dikirim ke Tempo dan divisualisasikan dalam bentuk Gantt Chart di Grafana.

## 📈 Monitoring Stack (Grafana)

Akses dashboard monitoring melalui:

- **URL**: `http://localhost:3000`
- **Username**: `admin`
- **Password**: `admin`

### Data Sources

1. **Prometheus**: Dashboard untuk grafik metrik.
2. **Loki**: Explore logs dengan query `{container="challenge_3-app-1"}`.
3. **Tempo**: Search traces berdasarkan `trace_id` atau Service Graph.

## 🛠️ Tech Stack

- **Go** (net/http + Go 1.22 enhanced routing)
- **PostgreSQL** — Database
- **Redis (go-redis/v9)** — Caching & Distributed Locks
- **OpenTelemetry** — Tracing SDK
- **Prometheus** — Metrics SDK
- **Loki & Promtail** — Log Management
- **Grafana** — Unified Visualization
- **Tempo** — Trace Storage
- **Docker & Docker Compose** — Containerization

## 📌 Fitur

- ✅ CRUD Rekening Bank
- ✅ Transfer antar rekening (dengan database transaction)
- ✅ Riwayat transaksi per akun
- ✅ Auto-migrate tabel saat startup
- ✅ Middleware Content-Type JSON
- ✅ Custom 404 handler
- ✅ Validasi input request
- ✅ Error handling yang konsisten
- ✅ Docker Containerization
- ✅ **Redis Caching**: Menyimpan respons akun dengan TTL 5 Menit di Memory
- ✅ **Redis Idempotency**: Mencegah request transfer ganda (Distributed Lock SetNX)
- ✅ **Global Timeout**: Membatasi maksimal waktu eksekusi API selama 10 Detik
- ✅ **Redis Rate Limiting**: Membatasi IP klien mengirim lebih dari 10 request dalam 5 detik
