package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func putExpenseHandler(c echo.Context) error {
	id := c.Param("id")

	e := Expense{}
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err_update := db.Prepare("UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1 RETURNING id")
	if err_update != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare update expense statment:" + err.Error()})
	}

	row := stmt.QueryRow(id, e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	err_scan := row.Scan(&e.ID)

	switch err_scan {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}