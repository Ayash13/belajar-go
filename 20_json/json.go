package jsoncodec

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Age      int      `json:"age,omitempty"`
	Password string   `json:"-"`
	Tags     []string `json:"tags,omitempty"`
}

func Run() {
	fmt.Println("=== JSON Encoding/Decoding ===")

	// ── Struct to JSON (Marshal) ──
	fmt.Println("\n--- Marshal (Struct → JSON) ---")
	user := User{
		ID:       1,
		Name:     "Ayash",
		Email:    "ayash@mail.com",
		Age:      22,
		Password: "secret123",
		Tags:     []string{"developer", "gopher"},
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("json.Marshal:", string(jsonBytes))
	fmt.Println("→ Password tidak muncul karena tag `json:\"-\"`")

	// ── Pretty print ──
	fmt.Println("\n--- MarshalIndent (Pretty) ---")
	pretty, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println(string(pretty))

	// ── JSON to Struct (Unmarshal) ──
	fmt.Println("\n--- Unmarshal (JSON → Struct) ---")
	jsonStr := `{"id":2,"name":"Budi","email":"budi@mail.com","age":25,"tags":["backend"]}`

	var decoded User
	err = json.Unmarshal([]byte(jsonStr), &decoded)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Decoded: %+v\n", decoded)

	// ── JSON to map (dynamic) ──
	fmt.Println("\n--- Unmarshal to map (Dynamic) ---")
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &result)
	fmt.Println("Name:", result["name"])
	fmt.Println("Age:", result["age"])

	// ── omitempty ──
	fmt.Println("\n--- omitempty ---")
	emptyUser := User{ID: 3, Name: "Citra", Email: "citra@mail.com"}
	emptyJSON, _ := json.Marshal(emptyUser)
	fmt.Println(string(emptyJSON))
	fmt.Println("→ Age dan Tags tidak muncul karena zero value + omitempty")

	// ── Slice of structs ──
	fmt.Println("\n--- Slice of Structs ---")
	users := []User{
		{ID: 1, Name: "Ayash", Email: "ayash@mail.com"},
		{ID: 2, Name: "Budi", Email: "budi@mail.com"},
	}
	usersJSON, _ := json.Marshal(users)
	fmt.Println(string(usersJSON))

	// ── Struct tags explained ──
	fmt.Println("\n=== Struct Tags ===")
	fmt.Println(`  json:"name"          → field name di JSON = "name"`)
	fmt.Println(`  json:"name,omitempty" → hilangkan jika zero value`)
	fmt.Println(`  json:"-"             → selalu skip (jangan serialize)`)
	fmt.Println(`  json:",string"       → encode angka sebagai string`)

	// ── In HTTP context ──
	fmt.Println("\n=== Di HTTP Handler ===")
	fmt.Println(`  // Decode request body`)
	fmt.Println(`  var user User`)
	fmt.Println(`  json.NewDecoder(r.Body).Decode(&user)`)
	fmt.Println(``)
	fmt.Println(`  // Encode response`)
	fmt.Println(`  w.Header().Set("Content-Type", "application/json")`)
	fmt.Println(`  json.NewEncoder(w).Encode(user)`)
}
