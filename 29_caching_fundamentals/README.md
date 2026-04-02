# Caching Fundamentals

Caching adalah teknik menyimpan salinan data yang sering diakses di lokasi penyimpanan sementara (cache) yang lebih cepat daripada sumber aslinya (seperti database).

## Mengapa Caching Penting?
1. **Meningkatkan Performa**: Mempercepat waktu respons aplikasi secara signifikan.
2. **Mengurangi Beban**: Mengurangi request berulang ke database atau layanan eksternal yang lambat/mahal.
3. **Skalabilitas**: Aplikasi dapat menangani lebih banyak trafik concurrent.

Contoh di `main.go` ini mendemonstrasikan bagaimana kita bisa membuat in-memory cache sederhana di level aplikasi menggunakan Go Maps dan `sync.RWMutex` untuk mencegah race conditions.
