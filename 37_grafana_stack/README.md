# Grafana Stack (Loki, Tempo, Prometheus, Grafana)

Grafana Stack adalah kumpulan tools observability dari Grafana Labs yang saling terintegrasi untuk meng-cover ketiga pilar observability: **Logs (Loki)**, **Traces (Tempo)**, dan **Metrics (Prometheus)**, semuanya divisualisasikan di satu dashboard **Grafana**.

## Arsitektur

```
┌──────────────────────────────────────────────────────────────────┐
│                        Go Application                            │
│  ┌────────────┐   ┌────────────────┐   ┌──────────────────────┐  │
│  │ Zap Logger │   │ OpenTelemetry  │   │ Prometheus Metrics   │  │
│  │ (stdout)   │   │ SDK (Traces)   │   │ (/metrics endpoint)  │  │
│  └─────┬──────┘   └──────┬─────────┘   └──────────┬───────────┘  │
└────────┼─────────────────┼────────────────────────┼──────────────┘
         │                 │                        │
         ▼                 ▼                        ▼
   ┌──────────┐     ┌──────────┐             ┌──────────┐
   │  Loki    │     │  Tempo   │             │Prometheus│
   │  :3100   │     │  :4318   │             │  :9090   │
   │  (Logs)  │     │ (Traces) │             │(Metrics) │
   └────┬─────┘     └────┬─────┘             └────┬─────┘
        │                │                        │
        └────────────────┼────────────────────────┘
                         ▼
                  ┌──────────────┐
                  │   Grafana    │
                  │   :3000      │
                  │ (Dashboard)  │
                  └──────────────┘
```

## Komponen

### 1. Loki — Log Aggregation
Loki mengumpulkan dan mengindeks log dari aplikasi. Tidak seperti Elasticsearch yang full-text index, Loki hanya mengindeks **label** (seperti `service_name`, `level`) sehingga jauh lebih ringan dan murah.

**Query Language**: LogQL
```logql
{service_name="bank-api"} |= "error"
{service_name="bank-api"} | json | level="error" | line_format "{{.message}}"
```

### 2. Tempo — Distributed Tracing Backend
Tempo menyimpan trace data yang dikirim via OpenTelemetry (OTLP). Tempo juga sangat ringan karena hanya membutuhkan object storage (atau local disk) tanpa indexing yang berat.

### 3. Prometheus — Metrics Collection
Prometheus melakukan scraping ke endpoint `/metrics` aplikasi secara berkala dan menyimpan data time-series. Query menggunakan PromQL.

### 4. Grafana — Unified Dashboard
Grafana menghubungkan ketiga data source di atas ke dalam satu UI. Fitur unggulannya:
- **Correlate**: Klik trace ID di log → langsung buka trace di Tempo
- **Explore**: Query real-time ke Loki, Tempo, dan Prometheus
- **Alerting**: Kirim alert ke Slack/Email jika metrics melebihi threshold

## Cara Menjalankan

```bash
# Jalankan seluruh stack
docker compose up -d

# Jalankan aplikasi Go
go run main.go

# Buat beberapa request untuk generate data
curl http://localhost:8080/accounts
curl http://localhost:8080/accounts/acc-001
curl -X POST http://localhost:8080/transfer -d '{"amount":100000}'
```

Setelah berjalan, akses Grafana di **http://localhost:3000** (login: `admin`/`admin`).

### Data Sources (sudah otomatis terkonfigurasi)
| Source | URL Internal | Fungsi |
|--------|-------------|--------|
| Loki | http://loki:3100 | Query logs |
| Tempo | http://tempo:3200 | Query traces |
| Prometheus | http://prometheus:9090 | Query metrics |

## Contoh Workflow Debugging di Grafana

1. **Alert**: Grafana mendeteksi error rate > 5% dari Prometheus
2. **Metrics**: Buka dashboard, lihat spike di `bank_api_http_requests_total{status="error"}`
3. **Logs**: Klik "Explore Loki" → filter `{service_name="bank-api"} | json | level="error"`
4. **Traces**: Dari log, klik trace_id → Grafana membuka trace di Tempo → terlihat span mana yang error

Contoh di `main.go` mendemonstrasikan aplikasi Go yang terintegrasi dengan ketiga pilar: Zap (→ Loki), OpenTelemetry (→ Tempo), dan Prometheus metrics, semuanya divisualisasikan di Grafana.
