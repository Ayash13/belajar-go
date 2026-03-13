package main

import (
	"fmt"
	"os"

	"belajar-go/01_variables"
	"belajar-go/02_constants"
	"belajar-go/03_functions"
	"belajar-go/04_conditions"
	"belajar-go/05_looping"
	"belajar-go/06_errors"
	"belajar-go/07_structs"
	"belajar-go/08_methods"
	"belajar-go/09_pointers"
	"belajar-go/10_interfaces"
	"belajar-go/11_dependency_injection"
	"belajar-go/practice_01_api_fetch"
)

var modules = map[string]struct {
	title string
	run   func()
}{
	"1":  {"Variables & Data Types", variables.Run},
	"2":  {"Constants, Iota & Operators", constants.Run},
	"3":  {"Functions", functions.Run},
	"4":  {"If/Else & Switch", conditions.Run},
	"5":  {"Looping", looping.Run},
	"6":  {"Basic Error Handling", errors.Run},
	"7":  {"Structs & Object Modeling", structs.Run},
	"8":  {"Methods (Value vs Pointer)", methods.Run},
	"9":  {"Pointer Concepts", pointers.Run},
	"10": {"Interfaces as Contracts", interfaces.Run},
	"11": {"Dependency Injection", dependency_injection.Run},
	"12": {"Practice 1: Simple API Fetch", practice_01_api_fetch.Run},
}

func main() {
	if len(os.Args) > 1 {
		key := os.Args[1]
		m, ok := modules[key]
		if !ok {
			fmt.Println("Usage: go run main.go [1-12]")
			fmt.Println("  1 - Variables    7 - Structs")
			fmt.Println("  2 - Constants    8 - Methods")
			fmt.Println("  3 - Functions    9 - Pointers")
			fmt.Println("  4 - Conditions  10 - Interfaces")
			fmt.Println("  5 - Looping     11 - Dependency Injection")
			fmt.Println("  6 - Errors      12 - Practice 1: API Fetch")
			return
		}
		fmt.Printf("--- %s. %s ---\n", key, m.title)
		m.run()
		return
	}

	fmt.Println("========================================")
	fmt.Println("  BELAJAR GO")
	fmt.Println("========================================")
	for i := 1; i <= 12; i++ {
		key := fmt.Sprintf("%d", i)
		m := modules[key]
		fmt.Printf("\n--- %s. %s ---\n", key, m.title)
		m.run()
	}
	fmt.Println("\n========================================")
	fmt.Println("  ALL COMPLETE!")
	fmt.Println("========================================")
}
