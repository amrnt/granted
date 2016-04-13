package main

import (
	"net/http"
	"strconv"

	"github.com/amrnt/granted"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

// User ...
type User struct {
	ID    int    `json:"id,omitempty" grant:"self"`
	Name  string `json:"name,omitempty" grant:"ALL"`
	Email string `json:"email,omitempty" grant:"self,ALL"`
	Phone string `json:"phone,omitempty" grant:"self"`
}

var users = []User{
	{
		ID:    1,
		Name:  "Amr",
		Email: "amr@email",
		Phone: "123123123",
	}, {
		ID:    2,
		Name:  "Julia",
		Email: "julia@email",
		Phone: "123123123",
	}, {
		ID:    3,
		Name:  "Ahmad",
		Email: "ahmad@email",
		Phone: "123123123",
	},
}

var authorize *granted.Authorize

func main() {
	e := echo.New()

	config := &granted.Config{}
	authorize = granted.New(config)

	// JWT or any token middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get JWT token and validate it
			// get current user id
			currentUserID := c.QueryParam("granted_to")
			c.Set("user_id", currentUserID)

			return next(c)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authorize.AddRole(granted.Role{
				Name: "self",
				CheckFunc: func(i granted.Instance) bool {
					return c.Get("user_id") == strconv.Itoa(i.(User).ID)
				},
			})

			return next(c)
		}
	})

	e.Get("/users", func(c echo.Context) error {
		return c.JSON(http.StatusOK, authorize.FilterToInterface(&users))
	})

	e.Get("/users/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		u := users[id-1]
		return c.JSON(http.StatusOK, authorize.FilterToInterface(&u))
	})

	e.Run(standard.New(":1323"))
}
