package main

// Currency Conversion

import "fmt"

// Struct untuk menampung hasil konversi
type ConversionResult struct {
	Currency string
	Amount   float64
}

// Mapping konversi
var rates = map[string]float64{
	"USD": 15000,
	"EUR": 16000,
	"JPY": 140,
	"SGD": 11000,
}

func main() {
	// Input nominal
	fmt.Println("Masukkan nominal:")
	var nominal float64
	fmt.Scanf("%f", &nominal)

	// Slice untuk simpan hasil konversi sebagai elemen dari struct
	results := []ConversionResult{}

	// Loop untuk konversi setiap mata uang
	for currency, rate := range rates {
		result := ConversionResult{
			Currency: currency,
			Amount:   nominal / rate,
		}
		results = append(results, result)
	}

	// Hasil konversi
	fmt.Println("===================")
	fmt.Println("Hasil Konversi:")
	for _, result := range results {
		fmt.Printf("%s: %.2f\n", result.Currency, result.Amount)
	}

}
