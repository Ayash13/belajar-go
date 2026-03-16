package httpserver

import "fmt"

func Run() {
	fmt.Println("=== Basic HTTP Server ===")

	fmt.Println("\n--- Membuat Server Sederhana ---")
	fmt.Println(`  package main`)
	fmt.Println(``)
	fmt.Println(`  import (`)
	fmt.Println(`      "fmt"`)
	fmt.Println(`      "net/http"`)
	fmt.Println(`  )`)
	fmt.Println(``)
	fmt.Println(`  func main() {`)
	fmt.Println(`      http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {`)
	fmt.Println(`          fmt.Fprintln(w, "Hello, World!")`)
	fmt.Println(`      })`)
	fmt.Println(``)
	fmt.Println(`      fmt.Println("Server running on :8080")`)
	fmt.Println(`      http.ListenAndServe(":8080", nil)`)
	fmt.Println(`  }`)

	fmt.Println("\n--- Komponen Utama ---")
	fmt.Println("  1. http.HandleFunc(pattern, handler) → Daftarkan route")
	fmt.Println("  2. http.ResponseWriter → Tulis response ke client")
	fmt.Println("  3. *http.Request       → Data request dari client")
	fmt.Println("  4. http.ListenAndServe → Start server pada port tertentu")

	fmt.Println("\n--- Default ServeMux ---")
	fmt.Println("  Go punya built-in router bernama DefaultServeMux.")
	fmt.Println("  http.HandleFunc() mendaftarkan handler ke DefaultServeMux.")
	fmt.Println("  Untuk proyek besar, gunakan custom mux atau library (gorilla/mux, chi).")

	fmt.Println("\n--- Graceful Shutdown ---")
	fmt.Println(`  srv := &http.Server{Addr: ":8080"}`)
	fmt.Println(`  go srv.ListenAndServe()`)
	fmt.Println(`  // ... wait for signal ...`)
	fmt.Println(`  srv.Shutdown(ctx) // graceful stop`)
}
