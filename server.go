package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TikhampornSky/assessment/repos"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

// var expenses_tmp = []Expense{
// 	{
// 		ID:     1,
// 		Title:  "strawberry smoothie",
// 		Amount: 79,
// 		Note:   "night market promotion discount 10 bath",
// 		Tags:   []string{"food", "beverage"},
// 	},
// }

func TokenCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token != "November 10, 2009" {
			return c.JSON(http.StatusUnauthorized, repos.Err{Message: "Unauthorization"})
		}
		return next(c)
	}
}

func main() {
	repos.InitDB()

	e := echo.New()

	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == "November 10, 2009wrong_token" {
	// 		return false, nil
	// 	}
	// 	return true, nil
	// }))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/expenses", TokenCheck(repos.GetExpensesHandler))
	e.GET("/expenses/:id", TokenCheck(repos.GetExpenseHandler))
	e.POST("/expenses", TokenCheck(repos.CreateExpenseHandler))
	e.PUT("/expenses/:id", TokenCheck(repos.PutExpenseHandler))

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	log.Fatal(e.Start(":2565"))
	log.Println("Exist! Byr bye")
}
