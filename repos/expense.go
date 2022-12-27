package repos

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) CreateExpenseHandler(c echo.Context) error {
	fmt.Println("Calling CreateExpenseHandler")
	e := Expense{}
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id", e.Title, e.Amount, e.Note, pq.Array(e.Tags))

	err = row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}

func (h *handler) GetExpensesHandler(c echo.Context) error {
	fmt.Println("Calling GetExpensesHandler")
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all expense statment:" + err.Error()})
	}
	fmt.Println("Calling GetExpensesHandler line 3")

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all expense:" + err.Error()})
	}

	expenses := []Expense{}
	for rows.Next() {
		e := Expense{}
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
		}
		expenses = append(expenses, e)
	}
	return c.JSON(http.StatusOK, expenses)
}

func (h *handler) GetExpenseHandler(c echo.Context) error {
	fmt.Println("Calling GetExpenseHandler")
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expense statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	e := Expense{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}

func (h *handler) PutExpenseHandler(c echo.Context) error {
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