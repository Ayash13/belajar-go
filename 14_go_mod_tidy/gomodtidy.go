package gomodtidy

import "fmt"

func Run() {
	fmt.Println("=== go mod tidy ===")

	fmt.Println("\n--- Apa itu go.mod? ---")
	fmt.Println("File go.mod mendefinisikan:")
	fmt.Println("  1. Nama module (misal: belajar-go)")
	fmt.Println("  2. Versi Go yang dipakai")
	fmt.Println("  3. Daftar dependency (library eksternal)")

	fmt.Println("\n--- Apa yang go mod tidy lakukan? ---")
	fmt.Println("  1. MENAMBAHKAN dependency yang dipakai tapi belum ada di go.mod")
	fmt.Println("  2. MENGHAPUS dependency yang di go.mod tapi sudah tidak dipakai")
	fmt.Println("  3. Mengupdate go.sum (checksum file)")

	fmt.Println("\n--- Contoh go.mod ---")
	fmt.Println(`  module belajar-go`)
	fmt.Println(``)
	fmt.Println(`  go 1.26.1`)
	fmt.Println(``)
	fmt.Println(`  require (`)
	fmt.Println(`      github.com/gorilla/mux v1.8.1`)
	fmt.Println(`      github.com/lib/pq v1.10.9`)
	fmt.Println(`  )`)

	fmt.Println("\n--- Perintah Penting ---")
	fmt.Println("  go mod init <module-name>  → Buat go.mod baru")
	fmt.Println("  go mod tidy               → Sinkronkan dependencies")
	fmt.Println("  go get <package>           → Tambah dependency baru")
	fmt.Println("  go mod download           → Download semua dependency")
	fmt.Println("  go mod verify             → Verifikasi checksum")

	fmt.Println("\n--- go.sum ---")
	fmt.Println("  File go.sum berisi hash checksum dari setiap dependency.")
	fmt.Println("  Ini memastikan dependency yang di-download tidak diubah/corrupted.")
	fmt.Println("  JANGAN edit go.sum secara manual!")

	fmt.Println("\n--- Kapan Jalankan go mod tidy? ---")
	fmt.Println("  1. Setelah menambah import package baru")
	fmt.Println("  2. Setelah menghapus import yang tidak dipakai")
	fmt.Println("  3. Sebelum commit ke git")
	fmt.Println("  4. Saat clone project baru")
}
