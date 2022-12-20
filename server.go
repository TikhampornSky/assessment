package main

import (
	"fmt"
	"log"
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

func main() {
	repos.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/expenses", repos.GetExpensesHandler)
	e.GET("/expenses/:id", repos.GetExpenseHandler)
	e.POST("/expenses", repos.CreateExpenseHandler)
	e.PUT("/expenses/:id", repos.PutExpenseHandler)

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	log.Fatal(e.Start(":2565"))
	log.Println("Exist! Byr bye")
}
