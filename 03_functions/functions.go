package functions

import "fmt"

// basic function
func greet(name string) string {
	return "Hello, " + name + "!"
}

// multiple return values
func divide(a, b float64) (float64, string) {
	if b == 0 {
		return 0, "cannot divide by zero"
	}
	return a / b, ""
}

// named return values
func rectangle(w, h float64) (area, perimeter float64) {
	area = w * h
	perimeter = 2 * (w + h)
	return
}

// variadic function
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func Run() {
	fmt.Println(greet("Ayash"))

	result, errMsg := divide(10, 3)
	fmt.Println("10 / 3 =", result, errMsg)

	_, errMsg2 := divide(10, 0)
	fmt.Println("10 / 0:", errMsg2)

	area, perimeter := rectangle(5, 3)
	fmt.Println("rectangle(5,3) -> area:", area, "perimeter:", perimeter)

	fmt.Println("sum(1,2,3,4,5) =", sum(1, 2, 3, 4, 5))

	// anonymous function
	double := func(x int) int { return x * 2 }
	fmt.Println("double(7) =", double(7))
}
