# Redis as In-Memory Data Store

Redis (Remote Dictionary Server) adalah in-memory data structure store yang open-source dan sering digunakan sebagai database, cache, dan message broker.

## Mengapa Redis?
Jika kita menggunakan *in-memory cache* local (seperti map di Go), setiap instance dari aplikasi kita akan memiliki cachenya masing-masing. Di sistem terdistribusi (microservices), ini akan menyebabkan data tidak sinkron. Redis menyediakan *centralized caching*, di mana setiap instance aplikasi menembak ke server Redis yang sama.

Program ini menggunakan package `github.com/redis/go-redis/v9` untuk berinteraksi dengan server Redis. Pastikan server Redis sudah menyala via Docker Compose sebelum menjalankannya:
```bash
# Menyalakan Redis
docker compose up -d

# Menjalankan aplikasi
go run main.go
```

## Parsing Struct ke dalam Cache
Pada contoh `main.go`, kita menyimpan `struct` kompleks (Nama, Umur, Pekerjaan) ke dalam Redis. 
Karena Redis menyimpan data dalam bentuk *bytes/string*, cara terbaik untuk menyimpan Object/Struct di Redis adalah mengonversinya menjadi bentuk **JSON** menggunakan `json.Marshal`. Ketika data ditarik kembali (Cache Hit), string tersebut akan di-*decode* ulang menggunakan `json.Unmarshal`.
