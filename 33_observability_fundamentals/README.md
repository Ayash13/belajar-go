# Observability Fundamentals (Logs, Metrics, Traces)

Observability adalah kemampuan untuk memahami kondisi internal dari sebuah sistem hanya dengan mengamati output-nya dari luar. Dalam dunia backend engineering, ada **3 pilar utama** observability yang dikenal sebagai **"Three Pillars of Observability"**.

## 3 Pilar Observability

### 1. Logs
Log adalah catatan detail dari setiap kejadian (event) dalam aplikasi — bisa berupa info, warning, atau error.

**Contoh**: "User ID 123 gagal login karena password salah pada 2026-04-01 09:00:00"

**Karakteristik:**
- Berbentuk teks atau JSON (structured log)
- Memberikan konteks **"apa yang terjadi"** secara spesifik
- Biasanya disimpan di Loki, Elasticsearch, atau CloudWatch

### 2. Metrics
Metrics adalah data numerik yang diukur dari waktu ke waktu — biasanya berupa angka yang bisa di-aggregate.

**Contoh**: "Total request per detik = 150", "CPU usage = 72%", "Rata-rata response time = 200ms"

**Karakteristik:**
- Berbentuk angka (counter, gauge, histogram)
- Memberikan gambaran **"seberapa sehat"** sistem secara keseluruhan
- Biasanya disimpan di Prometheus, InfluxDB, atau Datadog

### 3. Traces (Distributed Tracing)
Trace melacak perjalanan sebuah request dari awal hingga akhir melewati berbagai service/komponen.

**Contoh**: "Request GET /accounts membutuhkan 200ms total: 5ms di middleware → 15ms di service → 180ms di database"

**Karakteristik:**
- Terdiri dari `Trace ID` dan `Span` (setiap span = satu langkah operasi)
- Memberikan jawaban **"di mana bottleneck-nya"**
- Biasanya disimpan di Tempo, Jaeger, atau Zipkin

## Analogi Sederhana

Bayangkan aplikasi Anda adalah sebuah **Rumah Sakit**:

| Pilar | Analogi | Pertanyaan yang Dijawab |
|-------|---------|-------------------------|
| **Logs** | Rekam medis pasien | "Apa yang terjadi pada pasien ini?" |
| **Metrics** | Dashboard statistik RS | "Berapa total pasien hari ini? Berapa rata-rata waktu tunggu?" |
| **Traces** | Alur perjalanan pasien | "Pasien ini menghabiskan 2 jam di pendaftaran, 30 menit di lab, 1 jam di dokter" |

## Hubungan Antar Pilar

```
┌─────────────────────────────────────────────────────┐
│                   Request Masuk                      │
│                                                      │
│  ┌──────────┐   ┌──────────┐   ┌──────────────────┐ │
│  │  METRICS  │   │  TRACES  │   │      LOGS        │ │
│  │           │   │          │   │                   │ │
│  │ counter++ │   │ span_id: │   │ "Processing      │ │
│  │ latency:  │   │ abc-123  │   │  account ID 42"  │ │
│  │ 200ms     │   │          │   │                   │ │
│  └──────────┘   └──────────┘   └──────────────────┘ │
│                                                      │
│       Prometheus      Tempo/Jaeger     Loki/Zap      │
└─────────────────────────────────────────────────────┘
```

## Tools yang Akan Kita Pelajari

| Pilar | Tool | Fungsi |
|-------|------|--------|
| Logs | **Zap** (Go library) | Structured logging di aplikasi Go |
| Logs | **Loki** (Grafana) | Penyimpanan & query log terpusat |
| Metrics | **Prometheus** | Pengumpulan & penyimpanan metrics |
| Traces | **OpenTelemetry** | SDK untuk instrumentasi traces |
| Traces | **Tempo** (Grafana) | Penyimpanan & query distributed traces |
| Dashboard | **Grafana** | Visualisasi semua pilar di satu tempat |
| Testing | **k6** | Load testing & performance benchmarking |

## Menjalankan Contoh

```bash
go run main.go
```

Contoh di `main.go` mendemonstrasikan perbedaan antara _unstructured log_ (menggunakan `log`) dan _structured log_ (menggunakan `log/slog`), serta simulasi sederhana dari metric counter dan trace span.
