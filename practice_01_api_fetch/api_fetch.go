package practice_01_api_fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Todo struct matches the JSON response from jsonplaceholder
type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func Run() {
	fmt.Println("Fetching data from API...")

	// 1. Make the HTTP GET request
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	// Ensure the connection is closed after the function ends
	defer resp.Body.Close()

	// 2. Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// 3. Unmarshal (parse) the JSON into our struct
	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// 4. Print the result
	fmt.Println("\n--- Data Received ---")
	fmt.Printf("ID       : %d\n", todo.ID)
	fmt.Printf("User ID  : %d\n", todo.UserID)
	fmt.Printf("Title    : %s\n", todo.Title)
	fmt.Printf("Completed: %t\n", todo.Completed)
}
