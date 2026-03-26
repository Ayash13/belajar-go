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
	"belajar-go/12_packages"
	"belajar-go/13_exported"
	"belajar-go/14_go_mod_tidy"
	"belajar-go/15_separation_of_concerns"
	"belajar-go/16_error_wrapping"
	"belajar-go/17_database"
	"belajar-go/18_http_server"
	"belajar-go/19_handlers"
	"belajar-go/20_json"
	"belajar-go/21_http_methods"
	"belajar-go/22_status_codes"
	"belajar-go/23_middleware"
	"belajar-go/24_concurrency"
	"belajar-go/25_goroutines"
	"belajar-go/26_synchronization"
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
	"12": {"Package System", packages.Run},
	"13": {"Exported vs Unexported", exported.Run},
	"14": {"go mod tidy", gomodtidy.Run},
	"15": {"Separation of Concerns", separation.Run},
	"16": {"Error Wrapping", errorwrapping.Run},
	"17": {"Database Integration", database.Run},
	"18": {"Basic HTTP Server", httpserver.Run},
	"19": {"Handlers", handlers.Run},
	"20": {"JSON Encoding/Decoding", jsoncodec.Run},
	"21": {"HTTP Methods", httpmethods.Run},
	"22": {"Status Codes", statuscodes.Run},
	"23": {"Middleware Concepts", middleware.Run},
	"24": {"Concurrency vs Parallelism", concurrency.Run},
	"25": {"Go Routines", goroutines.Run},
	"26": {"Synchronization", synchronization.Run},
	"27": {"Practice 1: Simple API Fetch", practice_01_api_fetch.Run},
}

const totalModules = 27

func main() {
	if len(os.Args) > 1 {
		key := os.Args[1]
		m, ok := modules[key]
		if !ok {
			printUsage()
			return
		}
		fmt.Printf("--- %s. %s ---\n", key, m.title)
		m.run()
		return
	}

	fmt.Println("========================================")
	fmt.Println("  BELAJAR GO")
	fmt.Println("========================================")
	for i := 1; i <= totalModules; i++ {
		key := fmt.Sprintf("%d", i)
		m := modules[key]
		fmt.Printf("\n--- %s. %s ---\n", key, m.title)
		m.run()
	}
	fmt.Println("\n========================================")
	fmt.Println("  ALL COMPLETE!")
	fmt.Println("========================================")
}

func printUsage() {
	fmt.Printf("Usage: go run main.go [1-%d]\n", totalModules)
	fmt.Println("  1  - Variables          14 - go mod tidy")
	fmt.Println("  2  - Constants          15 - Separation of Concerns")
	fmt.Println("  3  - Functions          16 - Error Wrapping")
	fmt.Println("  4  - Conditions         17 - Database Integration")
	fmt.Println("  5  - Looping            18 - Basic HTTP Server")
	fmt.Println("  6  - Errors             19 - Handlers")
	fmt.Println("  7  - Structs            20 - JSON Encoding/Decoding")
	fmt.Println("  8  - Methods            21 - HTTP Methods")
	fmt.Println("  9  - Pointers           22 - Status Codes")
	fmt.Println("  10 - Interfaces         23 - Middleware")
	fmt.Println("  11 - Dependency Inject  24 - Concurrency vs Parallelism")
	fmt.Println("  12 - Package System     25 - Go Routines")
	fmt.Println("  13 - Exported           26 - Synchronization")
	fmt.Println("                          27 - Practice 1: API Fetch")
	fmt.Println("")
	fmt.Println("  Standalone Projects:")
	fmt.Println("    cd practice_02_postgres_crud && go run main.go")
	fmt.Println("    cd practice_03_nethttp_soc && go run main.go")
	fmt.Println("    cd challenge_1 && go run main.go")
	fmt.Println("    cd challenge_2 && go run main.go")
	fmt.Println("    cd challenge_3 && go run main.go")
	fmt.Println("    cd challenge_4 && go run .")
}
