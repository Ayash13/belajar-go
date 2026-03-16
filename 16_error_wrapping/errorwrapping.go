package errorwrapping

import (
	"errors"
	"fmt"
)

// ══════════════════════════════════════════════════
// Error Wrapping — Menambahkan konteks ke error
// ══════════════════════════════════════════════════

// ── Sentinel errors ──
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrValidation   = errors.New("validation error")
)

// ── Custom error type ──
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// ── Simulated layers ──

func findUserInDB(id int) (string, error) {
	if id != 1 {
		return "", ErrNotFound
	}
	return "Ayash", nil
}

func getUserFromRepo(id int) (string, error) {
	name, err := findUserInDB(id)
	if err != nil {
		// Wrap the error with context
		return "", fmt.Errorf("repo.GetUser(id=%d): %w", id, err)
	}
	return name, nil
}

func getUserFromService(id int) (string, error) {
	name, err := getUserFromRepo(id)
	if err != nil {
		// Wrap again with more context
		return "", fmt.Errorf("service.GetUser: %w", err)
	}
	return name, nil
}

func Run() {
	fmt.Println("=== Error Wrapping ===")

	// ── Success case ──
	fmt.Println("\n--- Case 1: User ditemukan ---")
	name, err := getUserFromService(1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("User:", name)
	}

	// ── Error case with wrapping ──
	fmt.Println("\n--- Case 2: User tidak ditemukan ---")
	_, err = getUserFromService(99)
	if err != nil {
		fmt.Println("Full error:", err)
		// Output: service.GetUser: repo.GetUser(id=99): not found

		// errors.Is — check if error chain contains a specific error
		fmt.Println("Is ErrNotFound?", errors.Is(err, ErrNotFound))
		fmt.Println("Is ErrUnauthorized?", errors.Is(err, ErrUnauthorized))
	}

	// ── Custom error type with Unwrap ──
	fmt.Println("\n--- Case 3: Custom AppError ---")
	appErr := &AppError{
		Code:    404,
		Message: "user not found",
		Err:     ErrNotFound,
	}
	fmt.Println("Error:", appErr)
	fmt.Println("Is ErrNotFound?", errors.Is(appErr, ErrNotFound))

	// errors.As — extract specific error type
	var target *AppError
	fmt.Println("Is AppError?", errors.As(appErr, &target))
	if target != nil {
		fmt.Println("  Code:", target.Code)
		fmt.Println("  Message:", target.Message)
	}

	// ── Wrapping with %w vs %v ──
	fmt.Println("\n--- Case 4: %w vs %v ---")

	baseErr := ErrValidation
	wrappedW := fmt.Errorf("handler: %w", baseErr) // preserves chain
	wrappedV := fmt.Errorf("handler: %v", baseErr) // breaks chain

	fmt.Println("With %%w — Is ErrValidation?", errors.Is(wrappedW, ErrValidation)) // true
	fmt.Println("With %%v — Is ErrValidation?", errors.Is(wrappedV, ErrValidation)) // false

	// ── Summary ──
	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. fmt.Errorf(\"context: %w\", err) → Wrap error dengan konteks")
	fmt.Println("2. errors.Is(err, target) → Cek apakah err chain mengandung target")
	fmt.Println("3. errors.As(err, &target) → Extract error type dari chain")
	fmt.Println("4. Gunakan %w (bukan %v) agar error chain tetap utuh")
	fmt.Println("5. Implement Unwrap() pada custom error type")
}
