package main

import (
	"fmt"

	"github.com/amrnt/granted"
)

// User ...
type User struct {
	ID    int    `grant:"ALL"`
	Name  string `grant:"ALL"`
	Email string `grant:"self"` // You can't see other's email (defined below)
	Admin bool   `grant:"NONE"`
}

func main() {
	currentUserID := 1

	granted.AddRole(granted.Role{
		Name: "self",
		CheckFunc: func(i granted.Instance) bool {
			return i.(User).ID == currentUserID
		},
	})

	//
	// Later on
	//

	u := User{
		ID:    1,
		Name:  "Amr",
		Email: "amr@email",
		Admin: true,
	}

	u2 := User{
		ID:    2,
		Name:  "Other user",
		Email: "other@email",
		Admin: false,
	}

	filtered := granted.Filter(&u).(*User)
	filtered2 := granted.Filter(&u2).(*User)

	//
	// Then
	//

	fmt.Println(filtered.ID == 1)              // true
	fmt.Println(filtered.Name == "Amr")        // true
	fmt.Println(filtered.Email == "amr@email") // true
	fmt.Println(filtered.Admin == false)       // true
	// It should return empty field for `Admin`.
	// But for now `Zero` of the value.
	// More: https://github.com/google/go-github/issues/19

	fmt.Println(filtered2.ID == 2)              // true
	fmt.Println(filtered2.Name == "Other user") // true
	fmt.Println(filtered2.Email == "")          // true
	fmt.Println(filtered2.Admin == false)       // true

	//
	// Fancy some JSON?
	//

	fmt.Println(granted.FilterToJSON(&u))
	fmt.Println(granted.FilterToJSON(&u2))
}
