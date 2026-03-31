# Docker dan Docker Compose

Materi ini membahas bagaimana membungkus (containerize) aplikasi Go menggunakan Docker dan menjalankannya bersama database PostgreSQL menggunakan Docker Compose.

## Komponen

1. **main.go**: Aplikasi Go HTTP Server sederhana yang terkoneksi ke PostgreSQL.
2. **Dockerfile**: File konfigurasi untuk membangun image Docker dari aplikasi Go.
3. **docker-compose.yml**: Konfigurasi untuk menjalankan multi-container (Aplikasi Go dan PostgreSQL).

## Cara Menjalankan

Buka terminal dan masuk ke direktori ini, kemudian jalankan perintah:

```bash
docker-compose up -d
```

Perintah di atas akan men-download image PostgreSQL, membangun (build) image aplikasi Go dari Dockerfile, dan menjalankan keduanya di background. Aplikasi akan menunggu sampai database siap menerima koneksi (jika sudah terhubung, akan mencetak log).

Setelah container berjalan, cek apakah aplikasi sudah terkoneksi ke database dengan mengakses endpoint:
```
http://localhost:8080/health
```

Jika sukses, endpoint akan me-return JSON:
```json
{"status":"ok","message":"Successfully connected to PostgreSQL"}
```

Untuk menghentikan dan menghapus container, jalankan:
```bash
docker-compose down
```
