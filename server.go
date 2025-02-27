package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TikhampornSky/assessment/repos"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	_ "github.com/lib/pq"
)

func main() {
	repos.InitDB()

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token != "November 10, 2009" {
				return c.JSON(http.StatusUnauthorized, "Unauthorized")
			}

			return next(c)
		}
	})

	fmt.Println("Set up server ok")

	e.GET("/expenses", repos.GetExpensesHandler)
	e.GET("/expenses/:id", repos.GetExpenseByIdHandler)
	e.POST("/expenses", repos.CreateExpenseHandler)
	e.PUT("/expenses/:id", repos.PutExpenseHandler)

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("Exist! Byr bye")
}
