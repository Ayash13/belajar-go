package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ── Model (GORM) ──

// Task represents the tasks table in the database
type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ── Database ──

var db *gorm.DB

func initDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=rahasiatau dbname=postgres sslmode=disable"
	}

	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate creates/updates the tasks table based on the Task struct
	err = db.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal("Failed to auto-migrate:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL (via GORM)")
}

// ── Handlers ──

func handleListTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	// Perform SELECT * FROM tasks
	if err := db.Order("id asc").Find(&tasks).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: tasks})
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "invalid id"})
		return
	}

	var task Task
	// Perform SELECT * FROM tasks WHERE id = ? LIMIT 1
	if err := db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Message: "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: task})
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "invalid JSON: " + err.Error()})
		return
	}
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "title is required"})
		return
	}

	task := Task{
		Title:       req.Title,
		Description: req.Description,
	}

	// Perform INSERT INTO tasks ...
	if err := db.Create(&task).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Message: "task created", Data: task})
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "invalid id"})
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "invalid JSON: " + err.Error()})
		return
	}

	// 1. Find existing task
	var task Task
	if err := db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Message: "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: err.Error()})
		return
	}

	// 2. Prepare updates
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Completed != nil {
		updates["completed"] = *req.Completed
	}

	// 3. Apply updates
	if len(updates) > 0 {
		if err := db.Model(&task).Updates(updates).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: err.Error()})
			return
		}
	}

	// Reload the task from DB to get the latest updated_at, etc.
	db.First(&task, id)

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Message: "task updated", Data: task})
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "invalid id"})
		return
	}

	// Perform DELETE FROM tasks WHERE id = ?
	result := db.Delete(&Task{}, id)
	if result.Error != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Message: "task not found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Message: "task deleted"})
}

// ── Middleware ──

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("%-6s %-20s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ── Helpers ──

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ── Main ──

func main() {
	initDB()

	// Get underlying sql.DB instance to defer close (GORM usually handles its own pooling)
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	mux := http.NewServeMux()

	// Register Routes
	mux.HandleFunc("GET /tasks", handleListTasks)
	mux.HandleFunc("GET /tasks/{id}", handleGetTask)
	mux.HandleFunc("POST /tasks", handleCreateTask)
	mux.HandleFunc("PUT /tasks/{id}", handleUpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handleDeleteTask)

	handler := corsMiddleware(loggingMiddleware(mux))

	fmt.Println("🚀 Server running on http://localhost:8080")
	fmt.Println("")
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /tasks      → List all tasks")
	fmt.Println("  GET    /tasks/{id} → Get task by ID")
	fmt.Println("  POST   /tasks      → Create task")
	fmt.Println("  PUT    /tasks/{id} → Update task")
	fmt.Println("  DELETE /tasks/{id} → Delete task")
	fmt.Println("")

	log.Fatal(http.ListenAndServe(":8080", handler))
}
