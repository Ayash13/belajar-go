package pointers

import "fmt"

func passByValue(val int) {
	val = 99 // only changes local copy
}

func passByPointer(ptr *int) {
	*ptr = 99 // changes the actual value
}

func Run() {
	x := 10
	fmt.Println("Initial x:", x)

	// memory address using &
	fmt.Println("Address of x:", &x)

	// pointer declaration
	var p *int = &x
	
	// dereferencing using *
	fmt.Println("Value at pointer p:", *p)

	// change value via pointer
	*p = 20
	fmt.Println("x after *p = 20:", x)

	// passing to functions
	num := 5
	passByValue(num)
	fmt.Println("\nAfter passByValue:", num) // still 5

	passByPointer(&num)
	fmt.Println("After passByPointer:", num) // changed to 99
}
