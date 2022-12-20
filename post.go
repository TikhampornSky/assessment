package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func createExpenseHandler(c echo.Context) error {
	var e Expense
	err := c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()} )
	}

	expenses = append(expenses, e)

	return c.JSON(http.StatusCreated, e)
}