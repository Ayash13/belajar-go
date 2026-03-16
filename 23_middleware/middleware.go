package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ══════════════════════════════════════════════════
// Middleware — Real Working Examples
// Starts a real HTTP server, applies middleware,
// and makes actual requests to demonstrate.
// ══════════════════════════════════════════════════

// ── Response helper ──

type jsonResponse struct {
	Message string `json:"message"`
	User    string `json:"user,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ══════════════════════════════════════════════════
// MIDDLEWARE 1: Logging
// ══════════════════════════════════════════════════

type logEntry struct {
	Method   string
	Path     string
	Status   int
	Duration time.Duration
}

var requestLogs []logEntry
var logMu sync.Mutex

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: 200}

		next.ServeHTTP(recorder, r)

		entry := logEntry{
			Method:   r.Method,
			Path:     r.URL.Path,
			Status:   recorder.status,
			Duration: time.Since(start),
		}
		logMu.Lock()
		requestLogs = append(requestLogs, entry)
		logMu.Unlock()
	})
}

// ══════════════════════════════════════════════════
// MIDDLEWARE 2: Auth (Token Check)
// ══════════════════════════════════════════════════

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			writeJSON(w, http.StatusUnauthorized, jsonResponse{Message: "missing Authorization header"})
			return // STOP — next is NOT called
		}

		if !strings.HasPrefix(token, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, jsonResponse{Message: "invalid token format, use: Bearer <token>"})
			return
		}

		actualToken := strings.TrimPrefix(token, "Bearer ")
		if actualToken != "secret-token-123" {
			writeJSON(w, http.StatusForbidden, jsonResponse{Message: "invalid token"})
			return
		}

		next.ServeHTTP(w, r) // authorized — continue
	})
}

// ══════════════════════════════════════════════════
// MIDDLEWARE 3: CORS
// ══════════════════════════════════════════════════

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return // preflight handled
		}

		next.ServeHTTP(w, r)
	})
}

// ══════════════════════════════════════════════════
// MIDDLEWARE 4: Rate Limiter (Simple)
// ══════════════════════════════════════════════════

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		// Clean old requests
		var valid []time.Time
		for _, t := range rl.requests[ip] {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		rl.requests[ip] = valid

		if len(valid) >= rl.limit {
			rl.mu.Unlock()
			writeJSON(w, http.StatusTooManyRequests, jsonResponse{
				Message: fmt.Sprintf("rate limit exceeded (%d requests per %v)", rl.limit, rl.window),
			})
			return
		}

		rl.requests[ip] = append(rl.requests[ip], now)
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// ══════════════════════════════════════════════════
// MIDDLEWARE 5: Recovery (Panic Handler)
// ══════════════════════════════════════════════════

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				writeJSON(w, http.StatusInternalServerError, jsonResponse{
					Message: fmt.Sprintf("internal server error: %v", err),
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ══════════════════════════════════════════════════
// Handlers
// ══════════════════════════════════════════════════

func publicHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, jsonResponse{Message: "this is a public endpoint"})
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, jsonResponse{Message: "welcome to protected area", User: "authenticated-user"})
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("something went terribly wrong!")
}

// ══════════════════════════════════════════════════
// Demo: make real HTTP requests
// ══════════════════════════════════════════════════

func makeRequest(method, url string, headers map[string]string) (int, string) {
	req, _ := http.NewRequest(method, url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, strings.TrimSpace(string(body))
}

func Run() {
	fmt.Println("=== Middleware — Real Working Examples ===")

	// Reset logs
	requestLogs = nil

	// ── Setup server with middleware chain ──
	limiter := newRateLimiter(3, 10*time.Second)

	mainMux := http.NewServeMux()

	// Public routes: CORS → Logging → Recovery → Handler
	publicChain := corsMiddleware(loggingMiddleware(recoveryMiddleware(http.HandlerFunc(publicHandler))))
	mainMux.Handle("GET /public", publicChain)
	mainMux.Handle("OPTIONS /public", publicChain) // for CORS preflight

	// Protected routes: CORS → Logging → Recovery → Auth → Handler
	protectedChain := corsMiddleware(loggingMiddleware(recoveryMiddleware(authMiddleware(http.HandlerFunc(protectedHandler)))))
	mainMux.Handle("GET /protected", protectedChain)

	// Panic route: Recovery catches the panic
	mainMux.Handle("GET /panic", loggingMiddleware(recoveryMiddleware(http.HandlerFunc(panicHandler))))

	// Rate limited route
	mainMux.Handle("GET /limited", loggingMiddleware(limiter.middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, jsonResponse{Message: "you got through!"})
	}))))

	// Start server
	server := &http.Server{Addr: ":19823", Handler: mainMux}
	go server.ListenAndServe()
	time.Sleep(100 * time.Millisecond) // wait for server to start

	base := "http://localhost:19823"

	// ── Test 1: Public endpoint ──
	fmt.Println("\n── 1. Logging Middleware ──")
	fmt.Println("   GET /public (no auth needed)")
	status, body := makeRequest("GET", base+"/public", nil)
	fmt.Printf("   Status: %d\n   Body: %s\n", status, body)

	// ── Test 2: Auth middleware — no token ──
	fmt.Println("\n── 2. Auth Middleware — No Token ──")
	fmt.Println("   GET /protected (tanpa Authorization header)")
	status, body = makeRequest("GET", base+"/protected", nil)
	fmt.Printf("   Status: %d\n   Body: %s\n", status, body)

	// ── Test 3: Auth middleware — wrong token ──
	fmt.Println("\n── 3. Auth Middleware — Wrong Token ──")
	fmt.Println("   GET /protected (token: Bearer wrong-token)")
	status, body = makeRequest("GET", base+"/protected", map[string]string{
		"Authorization": "Bearer wrong-token",
	})
	fmt.Printf("   Status: %d\n   Body: %s\n", status, body)

	// ── Test 4: Auth middleware — valid token ──
	fmt.Println("\n── 4. Auth Middleware — Valid Token ──")
	fmt.Println("   GET /protected (token: Bearer secret-token-123)")
	status, body = makeRequest("GET", base+"/protected", map[string]string{
		"Authorization": "Bearer secret-token-123",
	})
	fmt.Printf("   Status: %d\n   Body: %s\n", status, body)

	// ── Test 5: CORS headers ──
	fmt.Println("\n── 5. CORS Middleware ──")
	fmt.Println("   OPTIONS /public (preflight request)")
	req, _ := http.NewRequest("OPTIONS", base+"/public", nil)
	resp, _ := http.DefaultClient.Do(req)
	fmt.Printf("   Status: %d\n", resp.StatusCode)
	fmt.Printf("   Access-Control-Allow-Origin: %s\n", resp.Header.Get("Access-Control-Allow-Origin"))
	fmt.Printf("   Access-Control-Allow-Methods: %s\n", resp.Header.Get("Access-Control-Allow-Methods"))
	resp.Body.Close()

	// ── Test 6: Recovery middleware — catches panic ──
	fmt.Println("\n── 6. Recovery Middleware — Panic Handler ──")
	fmt.Println("   GET /panic (handler panics, but server survives)")
	status, body = makeRequest("GET", base+"/panic", nil)
	fmt.Printf("   Status: %d\n   Body: %s\n", status, body)

	// ── Test 7: Rate limiter ──
	fmt.Println("\n── 7. Rate Limiter Middleware ──")
	fmt.Println("   GET /limited (limit: 3 requests per 10s)")
	for i := 1; i <= 5; i++ {
		status, body = makeRequest("GET", base+"/limited", nil)
		fmt.Printf("   Request %d → %d: %s\n", i, status, body)
	}

	// ── Show collected logs ──
	fmt.Println("\n── Request Log (from Logging Middleware) ──")
	logMu.Lock()
	for _, l := range requestLogs {
		fmt.Printf("   %-6s %-15s → %d (%v)\n", l.Method, l.Path, l.Status, l.Duration.Round(time.Microsecond))
	}
	logMu.Unlock()

	// Shutdown server
	server.Close()

	// ── Summary ──
	fmt.Println("\n=== Middleware Pattern ===")
	fmt.Println("  Signature: func(http.Handler) http.Handler")
	fmt.Println("")
	fmt.Println("  Logging    → Catat method, path, status, duration")
	fmt.Println("  Auth       → Cek token, reject jika invalid")
	fmt.Println("  CORS       → Set Access-Control headers")
	fmt.Println("  Rate Limit → Batasi request per IP per waktu")
	fmt.Println("  Recovery   → Tangkap panic, return 500")
	fmt.Println("")
	fmt.Println("  Chain: cors(logging(recovery(auth(handler))))")
}
