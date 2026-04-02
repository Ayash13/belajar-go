# Time-To-Live (TTL)

Time-To-Live (TTL) adalah mekanisme pengaturan masa berlaku dari sebuah data di dalam cache. Ketika sebuah data mempunyai TTL 1 jam, maka otomatis setelah 1 jam berlalu, data tersebut akan hilang/dihapus oleh caching sistem.

Sangat penting untuk mencegah cache memakan memori habis (OOM), dan memastikan data menjadi fress/up-to-date (cache invalidation).

Aplikasi ini menggunakan Redis untuk mendemonstrasikan TTL.
