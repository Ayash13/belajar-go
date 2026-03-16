package httpmethods

import "fmt"

func Run() {
	fmt.Println("=== HTTP Methods (GET, POST, PUT, DELETE) ===")

	fmt.Println("\n--- HTTP Methods Overview ---")
	fmt.Println("  GET    → Ambil data (Read)")
	fmt.Println("  POST   → Buat data baru (Create)")
	fmt.Println("  PUT    → Update data keseluruhan (Update)")
	fmt.Println("  PATCH  → Update data sebagian (Partial Update)")
	fmt.Println("  DELETE → Hapus data (Delete)")

	fmt.Println("\n--- RESTful Routes ---")
	fmt.Println("  GET    /users      → List semua users")
	fmt.Println("  GET    /users/{id} → Get user by ID")
	fmt.Println("  POST   /users      → Create user baru")
	fmt.Println("  PUT    /users/{id} → Update user")
	fmt.Println("  DELETE /users/{id} → Delete user")

	fmt.Println("\n--- Go 1.22+ Routing ---")
	fmt.Println(`  mux := http.NewServeMux()`)
	fmt.Println(`  mux.HandleFunc("GET /users", listUsers)`)
	fmt.Println(`  mux.HandleFunc("GET /users/{id}", getUser)`)
	fmt.Println(`  mux.HandleFunc("POST /users", createUser)`)
	fmt.Println(`  mux.HandleFunc("PUT /users/{id}", updateUser)`)
	fmt.Println(`  mux.HandleFunc("DELETE /users/{id}", deleteUser)`)

	fmt.Println("\n--- Handler dengan Method Check (Pre Go 1.22) ---")
	fmt.Println(`  func usersHandler(w http.ResponseWriter, r *http.Request) {`)
	fmt.Println(`      switch r.Method {`)
	fmt.Println(`      case http.MethodGet:`)
	fmt.Println(`          // list users`)
	fmt.Println(`      case http.MethodPost:`)
	fmt.Println(`          // create user`)
	fmt.Println(`      default:`)
	fmt.Println(`          http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)`)
	fmt.Println(`      }`)
	fmt.Println(`  }`)

	fmt.Println("\n--- Path Parameters (Go 1.22+) ---")
	fmt.Println(`  func getUser(w http.ResponseWriter, r *http.Request) {`)
	fmt.Println(`      id := r.PathValue("id")`)
	fmt.Println(`      fmt.Fprintf(w, "User ID: %s", id)`)
	fmt.Println(`  }`)

	fmt.Println("\n--- Query Parameters ---")
	fmt.Println(`  func listUsers(w http.ResponseWriter, r *http.Request) {`)
	fmt.Println(`      page := r.URL.Query().Get("page")    // ?page=2`)
	fmt.Println(`      limit := r.URL.Query().Get("limit")  // ?limit=10`)
	fmt.Println(`  }`)

	fmt.Println("\n=== Idempotency ===")
	fmt.Println("  GET    → Idempotent ✅ (bisa dipanggil berulang, hasil sama)")
	fmt.Println("  PUT    → Idempotent ✅")
	fmt.Println("  DELETE → Idempotent ✅")
	fmt.Println("  POST   → NOT Idempotent ❌ (setiap call buat data baru)")
}
