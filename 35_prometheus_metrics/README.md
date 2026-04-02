# Prometheus Metrics

Prometheus adalah sistem monitoring open-source yang mengumpulkan metrics dari aplikasi dengan cara **scraping** endpoint `/metrics` secara berkala. Data yang sudah dikumpulkan bisa divisualisasikan menggunakan Grafana.

## Cara Kerja Prometheus

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Go App     │ scrape  │  Prometheus  │  query   │   Grafana    │
│ :8080/metrics├────────►│  :9090       ├────────►│   :3000      │
│              │  (pull) │              │ (PromQL) │              │
└──────────────┘         └──────────────┘         └──────────────┘
```

Prometheus menggunakan model **pull-based**: Prometheus yang aktif mengambil data dari aplikasi, bukan aplikasi yang mengirim data. Ini membuat arsitektur lebih simpel dan reliabel.

## 4 Tipe Metrics di Prometheus

| Tipe | Fungsi | Contoh |
|------|--------|--------|
| **Counter** | Angka yang hanya naik (monotonically increasing) | Total request, total error |
| **Gauge** | Angka yang bisa naik-turun | Jumlah goroutine aktif, CPU usage |
| **Histogram** | Distribusi nilai dalam bucket | Response time (p50, p95, p99) |
| **Summary** | Mirip histogram tapi menghitung quantile di client | Request duration percentiles |

## Menjalankan Contoh

Pastikan Docker sudah terinstal terlebih dahulu.

```bash
# Menjalankan Prometheus via Docker Compose
docker compose up -d

# Menjalankan aplikasi Go
go run main.go
```

Setelah aplikasi berjalan:
- **App metrics**: http://localhost:8080/metrics
- **Prometheus UI**: http://localhost:9090 (query metrics menggunakan PromQL)

## Contoh PromQL (Prometheus Query Language)

```promql
# Total request per detik
rate(bank_api_http_requests_total[1m])

# Rata-rata response time (p95)
histogram_quantile(0.95, rate(bank_api_http_duration_seconds_bucket[5m]))

# Error rate (%)
rate(bank_api_http_requests_total{status="error"}[5m]) / rate(bank_api_http_requests_total[5m]) * 100
```

Contoh di `main.go` mendemonstrasikan cara mengexpose metrics Counter, Gauge, dan Histogram dari sebuah HTTP server Go menggunakan library `prometheus/client_golang`.
