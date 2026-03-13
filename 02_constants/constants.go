package constants

import "fmt"

// iota auto-increments from 0
type Day int

const (
	Sunday Day = iota // 0
	Monday            // 1
	Tuesday           // 2
	Wednesday         // 3
)

const Pi = 3.14
const AppName = "BelajarGo"

func Run() {
	fmt.Println("Pi:", Pi)
	fmt.Println("AppName:", AppName)

	// iota
	fmt.Println("Sunday:", Sunday, "Monday:", Monday, "Tuesday:", Tuesday)

	// operators
	a, b := 10, 3
	fmt.Println("a + b =", a+b)
	fmt.Println("a - b =", a-b)
	fmt.Println("a * b =", a*b)
	fmt.Println("a / b =", a/b)
	fmt.Println("a % b =", a%b)

	// comparison
	fmt.Println("a > b:", a > b)
	fmt.Println("a == b:", a == b)

	// logical
	fmt.Println("true && false:", true && false)
	fmt.Println("true || false:", true || false)
	fmt.Println("!true:", !true)
}
