package conditions

import "fmt"

func Run() {
	score := 85

	// if/else
	if score >= 90 {
		fmt.Println("Grade: A")
	} else if score >= 80 {
		fmt.Println("Grade: B")
	} else if score >= 70 {
		fmt.Println("Grade: C")
	} else {
		fmt.Println("Grade: D")
	}

	// if with short statement
	if x := 10; x > 5 {
		fmt.Println(x, "is greater than 5")
	}

	// switch
	day := "Monday"
	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Println(day, "is a weekday")
	case "Saturday", "Sunday":
		fmt.Println(day, "is weekend")
	default:
		fmt.Println("unknown day")
	}

	// switch without expression (like if/else chain)
	temp := 35
	switch {
	case temp > 30:
		fmt.Println("hot!")
	case temp > 20:
		fmt.Println("warm")
	default:
		fmt.Println("cold")
	}
}
