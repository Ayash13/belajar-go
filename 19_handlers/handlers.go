package handlers

import "fmt"

func Run() {
	fmt.Println("=== HTTP Handlers ===")

	fmt.Println("\n--- 1. Function Handler ---")
	fmt.Println(`  http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {`)
	fmt.Println(`      fmt.Fprintln(w, "Hello!")`)
	fmt.Println(`  })`)

	fmt.Println("\n--- 2. Named Function Handler ---")
	fmt.Println(`  func helloHandler(w http.ResponseWriter, r *http.Request) {`)
	fmt.Println(`      fmt.Fprintln(w, "Hello!")`)
	fmt.Println(`  }`)
	fmt.Println(`  http.HandleFunc("/hello", helloHandler)`)

	fmt.Println("\n--- 3. Struct Handler (http.Handler interface) ---")
	fmt.Println(`  type GreetHandler struct {`)
	fmt.Println(`      Greeting string`)
	fmt.Println(`  }`)
	fmt.Println(``)
	fmt.Println(`  func (h *GreetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {`)
	fmt.Println(`      fmt.Fprintln(w, h.Greeting)`)
	fmt.Println(`  }`)
	fmt.Println(``)
	fmt.Println(`  http.Handle("/greet", &GreetHandler{Greeting: "Halo!"})`)

	fmt.Println("\n--- http.Handler Interface ---")
	fmt.Println(`  type Handler interface {`)
	fmt.Println(`      ServeHTTP(ResponseWriter, *Request)`)
	fmt.Println(`  }`)
	fmt.Println("  Semua HTTP handler di Go harus implement interface ini.")
	fmt.Println("  http.HandleFunc() hanyalah shortcut yang membungkus function biasa.")

	fmt.Println("\n--- Accessing Request Data ---")
	fmt.Println("  r.Method          → GET, POST, PUT, DELETE")
	fmt.Println("  r.URL.Path        → /users/123")
	fmt.Println("  r.URL.Query()     → query params (?name=Ayash)")
	fmt.Println("  r.Header.Get(key) → request headers")
	fmt.Println("  r.Body            → request body (io.ReadCloser)")
	fmt.Println("  r.PathValue(key)  → path params (Go 1.22+)")

	fmt.Println("\n--- Writing Response ---")
	fmt.Println("  w.Header().Set(key, value) → Set response header")
	fmt.Println("  w.WriteHeader(statusCode)  → Set status code")
	fmt.Println("  w.Write([]byte)            → Write response body")
	fmt.Println("  fmt.Fprintln(w, text)      → Write text response")
}
