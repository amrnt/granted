package main

import (
	"fmt"

	"github.com/amrnt/granted"
)

// User ...
type User struct {
	ID    uint   `grant:"nobody"`
	Name  string `grant:"everybody"`
	Email string `grant:"nobody"`
}

func main() {

	config := &granted.Config{}
	authorize := granted.New(config)

	authorize.AddRole(granted.Role{
		Name: "everybody",
		CheckFunc: func(i granted.Instance) bool {
			return true
		},
	})

	authorize.AddRole(granted.Role{
		Name: "nobody",
		CheckFunc: func(i granted.Instance) bool {
			return false
		},
	})

	authorize.AddRole(granted.Role{
		Name: "self",
		CheckFunc: func(i granted.Instance) bool {
			return i.(User).ID == 1
		},
	})

	// Later on

	u := User{
		ID:    1,
		Name:  "Amr",
		Email: "amr@email",
	}

	filtered := authorize.FilterToInterface(&u).(*User)

	// Then

	fmt.Println(filtered.ID == 1)              // false
	fmt.Println(filtered.Name == "Amr")        // true
	fmt.Println(filtered.Email == "amr@email") // false
}
