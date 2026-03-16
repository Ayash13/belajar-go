package exported

import "fmt"

// ── EXPORTED (Huruf kapital di awal) ──
// Bisa diakses dari package lain

// User is exported — accessible from outside this package.
type User struct {
	Name  string // exported field
	Email string // exported field
	age   int    // unexported field — only usable within this package
}

// NewUser is an exported constructor.
func NewUser(name, email string, age int) User {
	return User{
		Name:  name,
		Email: email,
		age:   age,
	}
}

// Greet is exported — callable from outside.
func (u User) Greet() string {
	return fmt.Sprintf("Hello, my name is %s (%s)", u.Name, u.Email)
}

// getAge is unexported — only callable within this package.
func (u User) getAge() int {
	return u.age
}

// GetAge is the exported way to access the unexported field.
func (u User) GetAge() int {
	return u.getAge()
}

// ── UNEXPORTED (Huruf kecil di awal) ──

// internalConfig is unexported — invisible outside this package.
type internalConfig struct {
	secret string
}

// newInternalConfig is unexported.
func newInternalConfig() internalConfig {
	return internalConfig{secret: "s3cr3t"}
}

// GetConfigSummary is an exported function that uses unexported internals.
func GetConfigSummary() string {
	cfg := newInternalConfig()
	return fmt.Sprintf("Config loaded (secret length: %d)", len(cfg.secret))
}

func Run() {
	fmt.Println("=== Exported vs Unexported ===")

	// Creating a User via exported constructor
	user := NewUser("Ayash", "ayash@mail.com", 22)

	// Accessing exported fields
	fmt.Println("Name:", user.Name)
	fmt.Println("Email:", user.Email)

	// user.age → COMPILE ERROR jika dari package lain
	// Tapi dari dalam package ini, bisa langsung akses:
	fmt.Println("Age (direct, same package):", user.age)

	// Via exported getter
	fmt.Println("Age (via GetAge):", user.GetAge())

	// Exported method
	fmt.Println(user.Greet())

	// Using unexported internals through exported function
	fmt.Println("\n" + GetConfigSummary())

	fmt.Println("\n=== Rules ===")
	fmt.Println("1. Huruf KAPITAL di awal → Exported (public)")
	fmt.Println("2. Huruf kecil di awal   → Unexported (private)")
	fmt.Println("3. Berlaku untuk: type, func, field, method, const, var")
	fmt.Println("4. Gunakan getter/setter untuk akses field unexported")
}
