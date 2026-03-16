package packages

import (
	"fmt"

	"belajar-go/12_packages/mathutil"
	"belajar-go/12_packages/stringutil"
)

func Run() {
	// Go organizes code into packages
	// Each folder = one package
	// Package name is declared at the top of every .go file

	// ── Using our custom packages ──
	fmt.Println("=== Custom Packages ===")

	sum := mathutil.Add(10, 5)
	fmt.Println("mathutil.Add(10, 5) =", sum)

	diff := mathutil.Subtract(10, 5)
	fmt.Println("mathutil.Subtract(10, 5) =", diff)

	avg := mathutil.Average([]float64{80, 90, 75, 88})
	fmt.Printf("mathutil.Average([80,90,75,88]) = %.2f\n", avg)

	// ── String utilities ──
	fmt.Println("\n=== String Utilities ===")

	reversed := stringutil.Reverse("Golang")
	fmt.Println("stringutil.Reverse(\"Golang\") =", reversed)

	capitalized := stringutil.Capitalize("hello world")
	fmt.Println("stringutil.Capitalize(\"hello world\") =", capitalized)

	wordCount := stringutil.WordCount("Belajar Go itu menyenangkan")
	fmt.Println("stringutil.WordCount(\"Belajar Go itu menyenangkan\") =", wordCount)

	// ── Key takeaways ──
	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. Satu folder = satu package")
	fmt.Println("2. Import menggunakan path: module/folder")
	fmt.Println("3. Nama package harus konsisten di semua file dalam folder")
	fmt.Println("4. Package 'main' adalah entry point program")
}
