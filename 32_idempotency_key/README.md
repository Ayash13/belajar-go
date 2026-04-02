# Idempotency Key & Data Deduplication

Idempotency adalah jaminan bahwa sebuah operasi dapat dijalankan berkali-kali namun memberikan hasil yang sama, dan state pada server tidak berubah jika dibandingkan dengan menjalankannya satu kali.

Ini vital di sistem pembayaran online, jika jaringan terputus (Network Timeout), klien akan meretry operasi HTTP ke server. Jika server tidak _idempotent_, maka klien bisa-bisa membayar dua kali untuk transaksi yang sama!

Redis sangat cocok untuk mencegah ini, memanfaatkan fungsi *SETNX* (Set if Not eXists) sebagai *Distributed Lock*.
