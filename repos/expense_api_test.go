// go:build expense

package repos

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func seedExpense(t *testing.T) Expense {
	var expense Expense
	body := bytes.NewBufferString(`{
		"title": "Sunflower Seed",
		"amount": 11.11,
		"note": "Pink Monday discount10%",
		"tags": ["food", "snack", "seed"]
	}`)

	rec := httptest.NewRecorder()
	c := setUpTest(http.MethodPost, uri("expenses"), body, rec)
	err := CreateExpenseHandler(c)

	if err != nil {
		t.Fatal("can't create expense!: ", err)
	}
	json.NewDecoder(rec.Body).Decode(&expense)

	return expense
}

func setUpTest(method, url string, body io.Reader, rec *httptest.ResponseRecorder) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(method, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009")
	c := e.NewContext(req, rec)

	return c
}

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie C++",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)

	rec := httptest.NewRecorder()
	c := setUpTest(http.MethodPost, uri("expenses"), body, rec)
	err := CreateExpenseHandler(c)

	var expense Expense
	json.NewDecoder(rec.Body).Decode(&expense)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, rec.Code)
	assert.NotEqual(t, 0, expense.ID)
	assert.Equal(t, "strawberry smoothie C++", expense.Title)
	assert.Equal(t, float64(79), expense.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", expense.Note)
	assert.Equal(t, []string{"food", "beverage"}, expense.Tags)
}

func TestGetAllExpenses(t *testing.T) {
	rec := httptest.NewRecorder()
	c := setUpTest(http.MethodGet, uri("expenses"), nil, rec)
	err := GetExpensesHandler(c)

	var expense []Expense
	json.NewDecoder(rec.Body).Decode(&expense)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, rec.Code)
	assert.Greater(t, len(expense), 0)
}

func TestGetExpenseByID(t *testing.T) {
	new_Expense := seedExpense(t)

	var latest Expense
	rec := httptest.NewRecorder()

	c := setUpTest(http.MethodGet, uri("expenses", strconv.Itoa(new_Expense.ID)), nil, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(new_Expense.ID))

	err := GetExpenseByIdHandler(c)

	json.NewDecoder(rec.Body).Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, new_Expense.ID, latest.ID)
	assert.Equal(t, "Sunflower Seed", latest.Title)
	assert.Equal(t, 11.11, latest.Amount)
	assert.Equal(t, "Pink Monday discount10%", latest.Note)
	assert.Equal(t, []string{"food", "snack", "seed"}, latest.Tags)
}

func TestUpdateExpenseByID(t *testing.T) {
	new_Expense := seedExpense(t)

	body := bytes.NewBufferString(`{
		"title": "Hamtaro Sunflower Seed",
		"amount": 40.5,
		"note": "Pink Monday discount10%",
		"tags": ["food", "snack", "putThing"]
	}`)

	var latest Expense
	rec := httptest.NewRecorder()
	c := setUpTest(http.MethodPut, uri("expenses", strconv.Itoa(new_Expense.ID)), body, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(new_Expense.ID))

	err := PutExpenseHandler(c)
	json.NewDecoder(rec.Body).Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, new_Expense.ID, latest.ID)
	assert.Equal(t, "Hamtaro Sunflower Seed", latest.Title)
	assert.Equal(t, 40.5, latest.Amount)
	assert.Equal(t, "Pink Monday discount10%", latest.Note)
	assert.Equal(t, []string{"food", "snack", "putThing"}, latest.Tags)
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}
