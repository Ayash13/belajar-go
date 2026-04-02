# OpenTelemetry Tracing

OpenTelemetry (OTel) adalah framework observability open-source yang menyediakan SDK dan API standar untuk mengumpulkan **Traces**, **Metrics**, dan **Logs** dari aplikasi. OpenTelemetry adalah proyek gabungan dari OpenTracing dan OpenCensus yang dikelola oleh CNCF.

## Mengapa OpenTelemetry?

Sebelum OTel, setiap vendor punya SDK sendiri (Jaeger SDK, Zipkin SDK, Datadog SDK). Jika kita ingin pindah vendor, kode harus diubah total. OpenTelemetry menyediakan **satu SDK universal** yang bisa mengirim data ke vendor manapun.

```
┌──────────────────────────────────────────────────────────┐
│                      Go Application                       │
│  ┌──────────────────────────────────────────────────────┐ │
│  │           OpenTelemetry SDK (Instrumentasi)          │ │
│  └──────────────────────┬───────────────────────────────┘ │
└─────────────────────────┼────────────────────────────────┘
                          │ OTLP (OpenTelemetry Protocol)
                          ▼
               ┌─────────────────────┐
               │   OTel Collector    │  (opsional, bisa langsung)
               └──────┬──────┬───────┘
                      │      │
            ┌─────────┘      └─────────┐
            ▼                          ▼
     ┌─────────────┐           ┌─────────────┐
     │    Tempo     │           │   Jaeger    │
     │  (Grafana)   │           │  (CNCF)     │
     └─────────────┘           └─────────────┘
```

## Konsep Utama

### 1. Trace
Satu trace = satu perjalanan lengkap sebuah request. Setiap trace punya **Trace ID** unik.

### 2. Span
Satu span = satu unit kerja di dalam trace. Setiap span memiliki:
- **Span ID**: ID unik untuk span ini
- **Parent Span ID**: ID span induknya (untuk membentuk tree/hierarchy)
- **Attributes**: Metadata tambahan (key-value)
- **Status**: OK, Error, atau Unset

### 3. Context Propagation
Trace ID diteruskan antar-function/service melalui Go `context.Context` sehingga semua span dalam satu request terhubung.

## Menjalankan Contoh

```bash
# Jalankan Tempo sebagai trace backend
docker compose up -d

# Jalankan aplikasi
go run main.go
```

Setelah berjalan, lakukan beberapa request:
```bash
curl http://localhost:8080/accounts
curl http://localhost:8080/accounts/acc-001
```

Contoh di `main.go` mendemonstrasikan cara menginstrumentasi HTTP server Go dengan OpenTelemetry SDK. Setiap request akan menghasilkan trace yang bisa dilihat di Tempo/Jaeger.
