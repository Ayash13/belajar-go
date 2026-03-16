package statuscodes

import "fmt"

func Run() {
	fmt.Println("=== HTTP Status Codes ===")

	fmt.Println("\n--- 1xx Informational ---")
	fmt.Println("  100 Continue")
	fmt.Println("  101 Switching Protocols")

	fmt.Println("\n--- 2xx Success ---")
	fmt.Println("  200 OK                → Request berhasil")
	fmt.Println("  201 Created           → Resource berhasil dibuat (POST)")
	fmt.Println("  204 No Content        → Berhasil tapi tidak ada body (DELETE)")

	fmt.Println("\n--- 3xx Redirection ---")
	fmt.Println("  301 Moved Permanently → URL berubah permanent")
	fmt.Println("  302 Found             → Redirect sementara")
	fmt.Println("  304 Not Modified      → Cache masih valid")

	fmt.Println("\n--- 4xx Client Error ---")
	fmt.Println("  400 Bad Request       → Request tidak valid")
	fmt.Println("  401 Unauthorized      → Belum login / token expired")
	fmt.Println("  403 Forbidden         → Tidak punya akses")
	fmt.Println("  404 Not Found         → Resource tidak ditemukan")
	fmt.Println("  405 Method Not Allowed→ Method HTTP salah")
	fmt.Println("  409 Conflict          → Konflik data (duplicate)")
	fmt.Println("  422 Unprocessable     → Validasi gagal")

	fmt.Println("\n--- 5xx Server Error ---")
	fmt.Println("  500 Internal Server Error → Bug di server")
	fmt.Println("  502 Bad Gateway           → Upstream server error")
	fmt.Println("  503 Service Unavailable   → Server overload/maintenance")

	fmt.Println("\n--- Penggunaan di Go ---")
	fmt.Println(`  // Constants dari net/http`)
	fmt.Println(`  http.StatusOK                  // 200`)
	fmt.Println(`  http.StatusCreated             // 201`)
	fmt.Println(`  http.StatusNoContent           // 204`)
	fmt.Println(`  http.StatusBadRequest          // 400`)
	fmt.Println(`  http.StatusUnauthorized        // 401`)
	fmt.Println(`  http.StatusForbidden           // 403`)
	fmt.Println(`  http.StatusNotFound            // 404`)
	fmt.Println(`  http.StatusInternalServerError // 500`)

	fmt.Println("\n--- Contoh di Handler ---")
	fmt.Println(`  // Success`)
	fmt.Println(`  w.WriteHeader(http.StatusOK)`)
	fmt.Println(`  json.NewEncoder(w).Encode(data)`)
	fmt.Println(``)
	fmt.Println(`  // Created`)
	fmt.Println(`  w.WriteHeader(http.StatusCreated)`)
	fmt.Println(`  json.NewEncoder(w).Encode(newUser)`)
	fmt.Println(``)
	fmt.Println(`  // Error`)
	fmt.Println(`  http.Error(w, "not found", http.StatusNotFound)`)

	fmt.Println("\n=== Tips ===")
	fmt.Println("  1. Selalu gunakan constant (http.StatusOK), bukan angka (200)")
	fmt.Println("  2. WriteHeader() harus dipanggil SEBELUM Write()")
	fmt.Println("  3. WriteHeader() hanya bisa dipanggil SEKALI per response")
}
