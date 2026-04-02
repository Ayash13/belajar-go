# Load Testing with k6

k6 adalah tool load testing open-source dari Grafana Labs yang ditulis dalam Go, dengan test scripts ditulis menggunakan JavaScript. k6 dirancang untuk developer dan CI/CD pipeline — ringan, scriptable, dan terintegrasi dengan Grafana Stack.

## Mengapa k6?

| Fitur | JMeter | Locust | **k6** |
|-------|--------|--------|--------|
| Script Language | XML/GUI | Python | **JavaScript** |
| Resource Usage | Berat | Sedang | **Ringan** |
| CLI Friendly | ❌ | ✅ | ✅ |
| CI/CD Ready | Ribet | ✅ | ✅ |
| Grafana Integration | ❌ | ❌ | ✅ (native) |

## Konsep Utama

### 1. Virtual Users (VUs)
VU adalah user virtual yang menjalankan script secara paralel. Setiap VU menjalankan satu iterasi script dari awal sampai akhir, kemudian mengulang.

### 2. Scenarios
Scenario menentukan pola traffic yang dihasilkan. Ada dua tipe utama:

#### a. Constant VUs
Jumlah VU tetap selama durasi test.
```
VUs ──────────────────────── (constant)
     10 VUs selama 30 detik
```

#### b. Ramping VUs
Jumlah VU naik/turun secara bertahap (stages).
```
VUs
 50 │          ┌──────────┐
    │         /            \
 10 │────────/              \──────
    └──────────────────────────────→ time
      ramp up    sustain   ramp down
```

### 3. Performance Metrics

| Metric | Penjelasan | Target Ideal |
|--------|-----------|--------------|
| **http_req_duration** | Waktu respons per request | p95 < 500ms |
| **http_req_failed** | Persentase request gagal | < 1% |
| **http_reqs** | Total request per detik (throughput) | Tergantung kebutuhan |
| **iteration_duration** | Waktu total satu iterasi VU | Tergantung scenario |
| **vus** | Jumlah VU aktif saat itu | Sesuai konfigurasi |

### 4. Checks & Thresholds
- **Checks**: Validasi respons (status code, body content) — seperti assertion di unit test.
- **Thresholds**: Batas performa yang harus dipenuhi — jika gagal, k6 exit code = non-zero (berguna di CI/CD).

## Instalasi k6

```bash
# macOS
brew install k6

# Docker (tanpa install)
docker run --rm -i grafana/k6 run - < scripts/smoke_test.js
```

## Cara Menjalankan

```bash
# Pastikan target API sudah berjalan
# (contoh: challenge_3 di localhost:8080)

# Smoke Test — test ringan untuk memastikan API tidak error
k6 run scripts/smoke_test.js

# Load Test — test dengan ramping VUs
k6 run scripts/load_test.js

# Stress Test — test batas maksimum sistem
k6 run scripts/stress_test.js
```

## Tipe Testing

| Tipe | Tujuan | VU Pattern |
|------|--------|------------|
| **Smoke Test** | Memastikan API tidak error | 1-2 VU, durasi pendek |
| **Load Test** | Mengukur performa normal | Ramping 10-50 VU |
| **Stress Test** | Mencari batas maksimum | Ramping hingga 100+ VU |
| **Spike Test** | Mengukur respons terhadap lonjakan tiba-tiba | 0 → 100 VU mendadak |
| **Soak Test** | Menemukan memory leak / degradasi | 50 VU, durasi panjang |

Folder `scripts/` berisi contoh test scripts untuk masing-masing tipe, lengkap dengan checks dan thresholds.
