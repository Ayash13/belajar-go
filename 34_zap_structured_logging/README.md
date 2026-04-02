# Structured Logging with Zap

Zap adalah library logging berkinerja tinggi dari Uber untuk Go. Dibandingkan `log` bawaan Go, Zap menawarkan **structured logging** (output JSON), **level-based logging** (Debug, Info, Warn, Error), dan performa yang sangat cepat karena minim alokasi memori.

## Mengapa Zap?

| Fitur | `log` (bawaan Go) | `slog` (Go 1.21+) | **Zap** |
|-------|--------------------|--------------------|---------|
| Structured (JSON) | ❌ | ✅ | ✅ |
| Log Levels | ❌ | ✅ | ✅ |
| Performa | Sedang | Bagus | **Sangat Cepat** |
| Field Types | ❌ | ✅ | ✅ (strongly typed) |
| Production Ready | ❌ | ✅ | ✅ (dipakai Uber) |

## Konsep Utama Zap

### 1. Logger Types
- **`zap.Logger`**: Logger utama, sangat cepat, menggunakan strongly-typed fields (`zap.String`, `zap.Int`).
- **`zap.SugaredLogger`**: Versi lebih fleksibel, bisa pakai `fmt.Sprintf`-style, sedikit lebih lambat.

### 2. Log Levels
```
Debug → Info → Warn → Error → DPanic → Panic → Fatal
```
Di production, biasanya level minimum di-set ke `Info` agar log `Debug` tidak ikut tercetak.

### 3. Output Format
- **Development**: Berwarna, human-readable, ada stacktrace untuk error.
- **Production**: JSON murni, cocok di-ingest oleh Loki/Elasticsearch.

## Cara Menjalankan

```bash
go run main.go
```

Contoh di `main.go` mendemonstrasikan:
- Perbedaan antara Development Logger dan Production Logger
- Penggunaan strongly-typed fields (`zap.String`, `zap.Int`, `zap.Duration`)
- Simulasi logging untuk operasi perbankan (transfer, login, error)
- Penggunaan `SugaredLogger` sebagai alternatif yang lebih fleksibel

## Integrasi dengan Loki
Di materi selanjutnya, output JSON dari Zap akan dikirim ke **Grafana Loki** agar bisa di-query dan divisualisasikan di Grafana Dashboard. Pola yang umum adalah: Aplikasi Go (Zap) → stdout → **Promtail/Alloy** → **Loki** → **Grafana**.
