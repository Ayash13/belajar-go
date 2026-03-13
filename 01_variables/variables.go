package variables

import "fmt"

func Run() {
	// var keyword
	var name string = "Ayash"
	var age int = 25
	var isStudent bool = true

	// short declaration
	score := 95.5
	city := "Jakarta"

	// zero values (default)
	var defaultInt int
	var defaultString string
	var defaultBool bool

	// multiple declaration
	var a, b, c int = 1, 2, 3

	fmt.Println("name:", name)
	fmt.Println("age:", age)
	fmt.Println("isStudent:", isStudent)
	fmt.Println("score:", score)
	fmt.Println("city:", city)
	fmt.Println("zero values ->", "int:", defaultInt, "string:", defaultString, "bool:", defaultBool)
	fmt.Println("multiple:", a, b, c)

	// type inference
	fmt.Printf("name type: %T, score type: %T\n", name, score)
}
