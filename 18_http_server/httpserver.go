package httpserver

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func Run() {
	fmt.Println("=== Basic HTTP Server — Real Working Example ===")

	// 1. Create a custom ServeMux (router)
	mux := http.NewServeMux()

	// 2. Register Routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintln(w, "Hello, World! This is the home page.")
	})

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is a basic Go HTTP server running in the background.")
	})

	// 3. Configure the HTTP Server
	server := &http.Server{
		Addr:         ":18080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 4. Start the server in a goroutine so it doesn't block
	go func() {
		fmt.Println("\n🚀 Starting actual server on http://localhost:18080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait a moment for server to start up
	time.Sleep(100 * time.Millisecond)

	// 5. Make real HTTP requests to our server to prove it works
	fmt.Println("\n── Making real HTTP requests to the running server ──")

	// Helper function to make request and print response
	makeRequest := func(url string) {
		fmt.Printf("\nTarget: %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("   Error:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("   Status: %s\n", resp.Status)
		fmt.Printf("   Body  : %s", string(body))
	}

	makeRequest("http://localhost:18080/")
	makeRequest("http://localhost:18080/about")
	makeRequest("http://localhost:18080/notfound")

	// 6. Graceful Shutdown
	fmt.Println("\n── Shutting down the server gracefully ──")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
	} else {
		fmt.Println("Server shut down successfully.")
	}

	// 7. Summary
	fmt.Println("\n=== Komponen Utama HTTP Server di Go ===")
	fmt.Println("  1. http.NewServeMux()  → Membuat router/multiplexer (direkomendasikan dibanding DefaultServeMux)")
	fmt.Println("  2. mux.HandleFunc()    → Mendaftarkan function handler ke specific path")
	fmt.Println("  3. &http.Server{}      → Konfigurasi server (port, timeout, handler)")
	fmt.Println("  4. ListenAndServe()    → Memulai server (blocking call)")
	fmt.Println("  5. server.Shutdown()   → Menghentikan server dengan anggun (graceful shutdown)")
}
