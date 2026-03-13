package looping

import "fmt"

func Run() {
	// classic for
	fmt.Print("classic: ")
	for i := 0; i < 5; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// while-style
	fmt.Print("while: ")
	count := 0
	for count < 3 {
		fmt.Print(count, " ")
		count++
	}
	fmt.Println()

	// range over slice
	fruits := []string{"apple", "banana", "mango"}
	for i, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}

	// range over map
	ages := map[string]int{"Ayash": 25, "Budi": 30}
	for name, age := range ages {
		fmt.Printf("  %s is %d\n", name, age)
	}

	// break & continue
	fmt.Print("skip even: ")
	for i := 0; i < 10; i++ {
		if i == 7 {
			break
		}
		if i%2 == 0 {
			continue
		}
		fmt.Print(i, " ")
	}
	fmt.Println()
}
