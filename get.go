package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getExpenseHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, expenses)
}