package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ══════════════════════════════════════════════════
// Database Integration — sqlx & GORM with PostgreSQL
// ══════════════════════════════════════════════════

// ── Models ──

// For sqlx — uses `db` struct tags
type TaskSqlx struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Completed   bool      `db:"completed"`
	CreatedAt   time.Time `db:"created_at"`
}

// For GORM — uses gorm conventions
type TaskGorm struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Completed   bool   `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TaskGorm) TableName() string {
	return "tasks_gorm"
}

func getConnStr() string {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=belajar_go sslmode=disable"
	}
	return connStr
}

// ══════════════════════════════════════════════════
// PART 1: sqlx
// ══════════════════════════════════════════════════

func demoSqlx(connStr string) {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║        SQLX — Enhanced database/sql  ║")
	fmt.Println("╚══════════════════════════════════════╝")

	// ── Connect ──
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Println("  ❌ Failed to connect:", err)
		return
	}
	defer db.Close()
	fmt.Println("  ✅ Connected to PostgreSQL via sqlx")

	// ── Create table ──
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS tasks_sqlx (
			id          SERIAL PRIMARY KEY,
			title       VARCHAR(255) NOT NULL,
			description TEXT DEFAULT '',
			completed   BOOLEAN DEFAULT FALSE,
			created_at  TIMESTAMP DEFAULT NOW()
		)
	`)
	db.MustExec("TRUNCATE tasks_sqlx RESTART IDENTITY")
	fmt.Println("  ✅ Table tasks_sqlx ready")

	// ── CREATE — NamedExec ──
	fmt.Println("\n  --- CREATE (NamedExec) ---")
	result, err := db.NamedExec(
		`INSERT INTO tasks_sqlx (title, description) VALUES (:title, :description)`,
		map[string]interface{}{
			"title":       "Belajar sqlx",
			"description": "Coba semua fitur sqlx",
		},
	)
	if err != nil {
		log.Println("  Insert error:", err)
		return
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("  Inserted %d row(s)\n", rows)

	// Insert more data
	db.MustExec(`INSERT INTO tasks_sqlx (title, description) VALUES ($1, $2)`, "Belajar GORM", "Coba GORM ORM")
	db.MustExec(`INSERT INTO tasks_sqlx (title, description) VALUES ($1, $2)`, "Build REST API", "Pakai semua konsep")

	// ── READ — Get (single row) ──
	fmt.Println("\n  --- READ — Get (single row → struct) ---")
	var task TaskSqlx
	err = db.Get(&task, "SELECT * FROM tasks_sqlx WHERE id = $1", 1)
	if err != nil {
		log.Println("  Get error:", err)
		return
	}
	fmt.Printf("  Found: [%d] %s — %s\n", task.ID, task.Title, task.Description)

	// ── READ — Select (multiple rows) ──
	fmt.Println("\n  --- READ — Select (multiple rows → slice) ---")
	var tasks []TaskSqlx
	err = db.Select(&tasks, "SELECT * FROM tasks_sqlx ORDER BY id")
	if err != nil {
		log.Println("  Select error:", err)
		return
	}
	for _, t := range tasks {
		status := "⬜"
		if t.Completed {
			status = "✅"
		}
		fmt.Printf("  %s [%d] %s\n", status, t.ID, t.Title)
	}

	// ── UPDATE ──
	fmt.Println("\n  --- UPDATE ---")
	_, err = db.Exec("UPDATE tasks_sqlx SET completed = TRUE WHERE id = $1", 1)
	if err != nil {
		log.Println("  Update error:", err)
		return
	}
	db.Get(&task, "SELECT * FROM tasks_sqlx WHERE id = $1", 1)
	fmt.Printf("  Updated: [%d] %s — completed: %v\n", task.ID, task.Title, task.Completed)

	// ── DELETE ──
	fmt.Println("\n  --- DELETE ---")
	res, _ := db.Exec("DELETE FROM tasks_sqlx WHERE id = $1", 3)
	affected, _ := res.RowsAffected()
	fmt.Printf("  Deleted %d row(s)\n", affected)

	// ── IN query ──
	fmt.Println("\n  --- SELECT IN (sqlx.In) ---")
	ids := []int{1, 2}
	query, args, _ := sqlx.In("SELECT * FROM tasks_sqlx WHERE id IN (?)", ids)
	query = db.Rebind(query) // convert ? to $1, $2 for postgres
	var filtered []TaskSqlx
	db.Select(&filtered, query, args...)
	for _, t := range filtered {
		fmt.Printf("  [%d] %s\n", t.ID, t.Title)
	}

	// ── NamedQuery ──
	fmt.Println("\n  --- NamedQuery ---")
	nrows, _ := db.NamedQuery(
		"SELECT * FROM tasks_sqlx WHERE completed = :completed",
		map[string]interface{}{"completed": true},
	)
	defer nrows.Close()
	for nrows.Next() {
		var t TaskSqlx
		nrows.StructScan(&t)
		fmt.Printf("  Completed: [%d] %s\n", t.ID, t.Title)
	}

	fmt.Println("\n  sqlx Key Points:")
	fmt.Println("  • db.Get()      → single row → struct")
	fmt.Println("  • db.Select()   → multiple rows → slice")
	fmt.Println("  • db.NamedExec  → named parameters (:title)")
	fmt.Println("  • sqlx.In()     → WHERE IN query helper")
	fmt.Println("  • Struct tag: `db:\"column_name\"`")
}

// ══════════════════════════════════════════════════
// PART 2: GORM
// ══════════════════════════════════════════════════

func demoGorm(connStr string) {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║        GORM — Full-featured ORM      ║")
	fmt.Println("╚══════════════════════════════════════╝")

	// ── Connect ──
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Println("  ❌ Failed to connect:", err)
		return
	}
	fmt.Println("  ✅ Connected to PostgreSQL via GORM")

	// ── AutoMigrate — creates/updates table from struct ──
	db.AutoMigrate(&TaskGorm{})
	db.Exec("TRUNCATE tasks_gorm RESTART IDENTITY")
	fmt.Println("  ✅ AutoMigrate: tasks_gorm table ready")

	// ── CREATE ──
	fmt.Println("\n  --- CREATE ---")
	task1 := TaskGorm{Title: "Belajar GORM", Description: "ORM lengkap untuk Go"}
	db.Create(&task1)
	fmt.Printf("  Created: [%d] %s (auto-generated ID!)\n", task1.ID, task1.Title)

	// Batch create
	tasks := []TaskGorm{
		{Title: "Belajar Migrations", Description: "AutoMigrate dan manual migration"},
		{Title: "Belajar Associations", Description: "HasOne, HasMany, BelongsTo"},
	}
	db.Create(&tasks)
	fmt.Printf("  Batch created: %d tasks\n", len(tasks))

	// ── READ ──
	fmt.Println("\n  --- READ — First (single) ---")
	var found TaskGorm
	db.First(&found, 1) // find by primary key
	fmt.Printf("  Found: [%d] %s\n", found.ID, found.Title)

	fmt.Println("\n  --- READ — Find (all) ---")
	var allTasks []TaskGorm
	db.Find(&allTasks)
	for _, t := range allTasks {
		status := "⬜"
		if t.Completed {
			status = "✅"
		}
		fmt.Printf("  %s [%d] %s\n", status, t.ID, t.Title)
	}

	fmt.Println("\n  --- READ — Where (filtered) ---")
	var incomplete []TaskGorm
	db.Where("completed = ?", false).Find(&incomplete)
	fmt.Printf("  Incomplete tasks: %d\n", len(incomplete))

	// ── UPDATE ──
	fmt.Println("\n  --- UPDATE ---")
	// Update single field
	db.Model(&TaskGorm{}).Where("id = ?", 1).Update("completed", true)
	// Update multiple fields
	db.Model(&TaskGorm{}).Where("id = ?", 2).Updates(map[string]interface{}{
		"completed":   true,
		"description": "Updated description!",
	})

	db.Find(&allTasks)
	for _, t := range allTasks {
		status := "⬜"
		if t.Completed {
			status = "✅"
		}
		fmt.Printf("  %s [%d] %s — %s\n", status, t.ID, t.Title, t.Description)
	}

	// ── DELETE ──
	fmt.Println("\n  --- DELETE ---")
	db.Delete(&TaskGorm{}, 3) // delete by primary key
	var remaining []TaskGorm
	db.Find(&remaining)
	fmt.Printf("  Remaining: %d tasks\n", len(remaining))

	// ── Counting ──
	fmt.Println("\n  --- Count ---")
	var count int64
	db.Model(&TaskGorm{}).Where("completed = ?", true).Count(&count)
	fmt.Printf("  Completed tasks: %d\n", count)

	fmt.Println("\n  GORM Key Points:")
	fmt.Println("  • db.Create()      → INSERT")
	fmt.Println("  • db.First()       → SELECT ... LIMIT 1")
	fmt.Println("  • db.Find()        → SELECT *")
	fmt.Println("  • db.Where().Find()→ SELECT ... WHERE")
	fmt.Println("  • db.Update()      → UPDATE single field")
	fmt.Println("  • db.Updates()     → UPDATE multiple fields")
	fmt.Println("  • db.Delete()      → DELETE")
	fmt.Println("  • db.AutoMigrate() → Create/alter table from struct")
}

// ══════════════════════════════════════════════════
// Comparison
// ══════════════════════════════════════════════════

func printComparison() {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║      sqlx vs GORM — Comparison       ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println("  ┌──────────────────┬──────────────────┬──────────────────┐")
	fmt.Println("  │                  │ sqlx             │ GORM             │")
	fmt.Println("  ├──────────────────┼──────────────────┼──────────────────┤")
	fmt.Println("  │ Style            │ SQL-first        │ ORM (method)     │")
	fmt.Println("  │ Query            │ Raw SQL          │ Chain methods    │")
	fmt.Println("  │ Migration        │ Manual           │ AutoMigrate      │")
	fmt.Println("  │ Learning curve   │ Mudah            │ Sedang           │")
	fmt.Println("  │ Control          │ Tinggi           │ Sedang           │")
	fmt.Println("  │ Boilerplate      │ Lebih banyak     │ Sedikit          │")
	fmt.Println("  │ Performance      │ Lebih cepat      │ Sedikit overhead │")
	fmt.Println("  │ Cocok untuk      │ Simple projects  │ Complex projects │")
	fmt.Println("  └──────────────────┴──────────────────┴──────────────────┘")
	fmt.Println("")
	fmt.Println("  Pilih sqlx jika: suka nulis SQL sendiri, butuh kontrol penuh")
	fmt.Println("  Pilih GORM jika: mau cepat develop, butuh migration & associations")
}

func Run() {
	fmt.Println("=== Database Integration (PostgreSQL) ===")

	connStr := getConnStr()
	demoSqlx(connStr)
	demoGorm(connStr)
	printComparison()
}
