package errors

import (
	"errors"
	"fmt"
)

// function that returns an error
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// custom error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "cannot be negative"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "unrealistic value"}
	}
	return nil
}

func Run() {
	// basic error handling
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 3 =", result)
	}

	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// custom error
	if err := validateAge(-5); err != nil {
		fmt.Println("Validation error:", err)
	}

	if err := validateAge(200); err != nil {
		fmt.Println("Validation error:", err)
	}

	fmt.Println("validateAge(25):", validateAge(25))
}
