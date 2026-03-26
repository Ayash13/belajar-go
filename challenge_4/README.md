# Challenge 4: Banking Transaction Processing

## Latar Belakang

Sebuah sistem perbankan ingin memproses transaksi transfer antar rekening secara **parallel** (semua transaksi diproses bersamaan). Ini mensimulasikan kondisi nyata dimana banyak nasabah melakukan transfer di waktu yang sama.

## Aturan Sistem

1. **Fully Parallel** — Jumlah worker = jumlah transaksi. Setiap transaksi langsung diproses oleh worker-nya sendiri, tidak perlu mengantri.

2. **Timeout 4 detik** — Setiap transaksi punya batas waktu maksimal 4 detik. Jika proses melebihi batas ini, transaksi otomatis dianggap **TIMEOUT** dan dibatalkan.

3. **Waktu proses random 1-5 detik** — Karena komunikasi antar bank membutuhkan waktu yang berbeda-beda, setiap transaksi disimulasikan dengan waktu proses acak antara 1 sampai 5 detik. Artinya ada kemungkinan transaksi selesai tepat waktu, dan ada yang timeout.

4. **Validasi saldo** — Jika transaksi selesai dalam batas waktu, sistem mengecek apakah saldo pengirim cukup:
   - **Saldo cukup → SUCCESS**: saldo pengirim dikurangi, saldo penerima ditambah, tampilkan saldo terbaru kedua rekening
   - **Saldo tidak cukup → FAILED**: transaksi dibatalkan, tampilkan info gagal beserta saldo saat ini

5. **Data konsisten** — Karena semua transaksi berjalan bersamaan, akses ke data saldo dilindungi dengan `sync.Mutex` agar tidak terjadi race condition (data corrupt).

## Sample Data

### Rekening

| Bank | No Rekening | Saldo Awal |
|------|-------------|------------|
| CIMB | C001 | Rp 300.000 |
| MANDIRI | M002 | Rp 500.000 |
| BNI | B003 | Rp 400.000 |
| BCA | BC04 | Rp 800.000 |

**Total saldo awal: Rp 2.000.000** (harus tetap sama di akhir, kecuali ada timeout)

### Daftar Transaksi

| ID | Dari | Ke | Nominal |
|----|------|----|---------| 
| 1 | C001 | M002 | Rp 300.000 |
| 2 | C001 | B003 | Rp 600.000 |
| 3 | M002 | BC04 | Rp 200.000 |
| 4 | B003 | C001 | Rp 500.000 |
| 5 | BC04 | M002 | Rp 700.000 |
| 6 | C001 | M002 | Rp 400.000 |
| 7 | B003 | C001 | Rp 300.000 |

> **Catatan:** TX-2 (Rp 600.000 dari C001) kemungkinan besar akan FAILED karena saldo C001 hanya Rp 300.000. Tapi karena transaksi berjalan bersamaan, hasilnya bisa berbeda setiap kali dijalankan tergantung urutan eksekusi.

## Cara Jalankan

```bash
cd challenge_4
go run .
```

## Alur Program

```
main()
  │
  ├── Buat channel `jobs` dan `results` (buffered)
  │
  ├── Jalankan N worker goroutine (N = jumlah transaksi)
  │     ├── Worker 1 ─┐
  │     ├── Worker 2 ──┤
  │     ├── Worker 3 ──┤
  │     ├── ...        ├── Semua worker siap mengambil dari `jobs`
  │     ├── Worker 6 ──┤
  │     └── Worker 7 ─┘
  │
  ├── Kirim 7 transaksi ke channel `jobs` → langsung diambil worker
  │
  ├── Setiap worker memproses transaksi secara parallel:
  │     ├── Buat context.WithTimeout(4 detik)
  │     ├── Simulasi proses (random 1-5 detik)
  │     └── select {
  │           case selesai < 4 detik → lock mutex → cek saldo → SUCCESS/FAILED
  │           case timeout > 4 detik → TIMEOUT
  │         }
  │
  ├── Tunggu semua worker selesai (WaitGroup)
  │
  └── Tampilkan saldo akhir semua rekening + total
```

## Konsep yang Digunakan

| Konsep | Apa | Dipakai Untuk |
|--------|-----|---------------|
| **Goroutines** | Fungsi yang berjalan concurrent | N worker (1 per transaksi) memproses semua transaksi secara parallel |
| **Channels** | Pipa komunikasi antar goroutine | `jobs` = antrian transaksi, `results` = hasil proses |
| **Select** | Menunggu salah satu dari beberapa channel | Memilih antara "proses selesai" atau "timeout" |
| **Worker Pool** | Pola: N worker mengambil job dari antrian | N worker = N transaksi, semua berjalan bersamaan |
| **Context** | Kontrol lifecycle goroutine (cancel/timeout) | `WithTimeout(4s)` membatasi waktu setiap transaksi |
| **Mutex** | Kunci akses ke shared data | Melindungi saldo rekening agar tidak corrupt saat diakses bersamaan |
| **WaitGroup** | Menunggu kumpulan goroutine selesai | Menunggu semua worker selesai sebelum print saldo akhir |

## Contoh Output

```
Worker-2 | Mengambil TX-1
Worker-5 | Mengambil TX-5
Worker-3 | Mengambil TX-3
Worker-19 | Mengambil TX-2
Worker-4 | Mengambil TX-4
10:59:45 | Processing TX-1 | C001 -> M002 | Rp300000
Worker-7 | Mengambil TX-7
10:59:45 | Processing TX-3 | M002 -> BC04 | Rp200000
10:59:45 | Processing TX-2 | C001 -> B003 | Rp600000
10:59:45 | Processing TX-5 | BC04 -> M002 | Rp700000
10:59:45 | Processing TX-7 | B003 -> C001 | Rp300000
10:59:45 | Processing TX-4 | B003 -> C001 | Rp500000
Worker-6 | Mengambil TX-6
10:59:45 | Processing TX-6 | C001 -> M002 | Rp400000
10:59:46 | TX-4 FAILED | B003 -> C001 | Transfer Rp500000 | Saldo B003 hanya Rp400000
10:59:47 | TX-5 SUCCESS
    Saldo BC04 sekarang Rp100000
    Saldo M002 sekarang Rp1200000
10:59:47 | TX-7 SUCCESS
    Saldo B003 sekarang Rp100000
    Saldo C001 sekarang Rp600000
10:59:48 | TX-2 SUCCESS
    Saldo C001 sekarang Rp0
    Saldo B003 sekarang Rp700000
10:59:49 | TX-3 TIMEOUT
10:59:49 | TX-6 TIMEOUT
10:59:49 | TX-1 TIMEOUT

===== FINAL BALANCE =====
CIMB (C001) : Rp0
MANDIRI (M002) : Rp1200000
BNI (B003) : Rp700000
BCA (BC04) : Rp100000

Total Saldo Semua Bank : Rp2000000
```

## Kemungkinan Hasil

Karena waktu proses random dan semua transaksi berjalan parallel, setiap kali dijalankan hasilnya **bisa berbeda**:

- Semua 7 transaksi mulai diproses hampir bersamaan
- Transaksi yang selesai duluan bisa menghabiskan saldo lebih dulu
- Transaksi berikutnya dari rekening yang sama bisa gagal karena saldo sudah berkurang
- Transaksi dengan waktu proses > 4 detik akan selalu TIMEOUT
- Total saldo akhir bisa berubah jika ada transaksi yang TIMEOUT (uang "menghilang" karena tidak diproses)
