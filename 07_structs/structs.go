package structs

import "fmt"

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	IsActive  bool
}

type Order struct {
	OrderID string
	User    User // nested struct
	Total   float64
}

func Run() {
	// initialization
	user1 := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		IsActive:  true,
	}

	// zero values
	var user2 User
	user2.ID = 2
	user2.FirstName = "Jane"

	// nested struct
	order := Order{
		OrderID: "ORD-001",
		User:    user1,
		Total:   150.50,
	}

	fmt.Printf("User 1: %+v\n", user1)
	fmt.Printf("User 2: %+v\n", user2)
	fmt.Printf("Order : %+v\n", order)
	fmt.Println("Order's User Email:", order.User.Email)
}
