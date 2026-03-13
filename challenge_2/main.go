package main

// Card Identification
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Struct untuk menampung hasil konversi
type CardType struct {
	Prefix []string
	Length []int
	Name   []string
}

// Map jenis kartu, prefix, panjang nomor kartu
var cardTypes = map[string]CardType{
	"China UnionPay": {
		Prefix: []string{"62"},
		Length: []int{16, 17, 18, 19},
	},
	"Switch": {
		Prefix: []string{"4903", "4905", "4911", "4936", "564182", "633110", "6333", "6759"},
		Length: []int{16, 18, 19},
	},
}

// Input nomor kartu lebih dari satu
func main() {
	fmt.Println("Masukkan nomor kartu:")
	var input []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		input = append(input, strings.Fields(line)...)
	}

	// Validasi nomor kartu bedasarkan map
	for _, card := range input {
		found := false
		for cardType, cardTypeValue := range cardTypes {
			for _, prefix := range cardTypeValue.Prefix {
				if strings.HasPrefix(card, prefix) {
					// Periksa apakah panjang cocok dengan yang ada di map
					for _, validLen := range cardTypeValue.Length {
						if len(card) == validLen {
							fmt.Printf("Nomor %s adalah %s\n", card, cardType)
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			fmt.Printf("Jenis kartu tidak dikenali (nomor: %s)\n", card)
		}
	}

}
