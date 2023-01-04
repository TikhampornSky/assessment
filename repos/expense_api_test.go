// go:build unitAPI

package repos

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setUpTest(method, url string, body io.Reader, rec *httptest.ResponseRecorder) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(method, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009")
	c := e.NewContext(req, rec)

	return c
}

func setUpDB(t *testing.T) sqlmock.Sqlmock {
	db, mock, errMock := sqlmock.New()
	if errMock != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errMock)
	}
	SetDB(db)

	return mock
}

func TestCreateExpense(t *testing.T) {

	mock := setUpDB(t)

	var data = Expense{
		Title:  "strawberry smoothie C++",
		Amount: 79.0,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "snack"},
	}

	mockedSql := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)"
	mockedRow := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(mockedSql)).WithArgs(data.Title, data.Amount, data.Note, pq.Array(data.Tags)).WillReturnRows((mockedRow))

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie C++",
		"amount": 79.0,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "snack"]
	}`)

	rec := httptest.NewRecorder()
	c := setUpTest(http.MethodPost, uri("expenses"), body, rec)
	errCreate := CreateExpenseHandler(c)
	assert.NoError(t, errCreate)

	var expense Expense
	json.NewDecoder(rec.Body).Decode(&expense)

	assert.EqualValues(t, http.StatusCreated, rec.Code)
	assert.NotEqual(t, data.ID, expense.ID)
	assert.Equal(t, data.Title, expense.Title)
	assert.Equal(t, data.Amount, expense.Amount)
	assert.Equal(t, data.Note, expense.Note)
	assert.Equal(t, data.Tags, expense.Tags)
}

func TestGetAllExpenses(t *testing.T) {

	mock := setUpDB(t)

	expensesMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "test-title", 99.99, "test-note", pq.Array([]string{"food", "snack"}))
	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").ExpectQuery().WillReturnRows(expensesMockRows)

	rec := httptest.NewRecorder()
	c := setUpTest(http.MethodGet, uri("expenses"), nil, rec)
	err := GetExpensesHandler(c)

	expected := "[{\"id\":1,\"title\":\"test-title\",\"amount\":99.99,\"note\":\"test-note\",\"tags\":[\"food\",\"snack\"]}]"
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
}

func TestGetExpenseByID(t *testing.T) {

	mock := setUpDB(t)
	expensesMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "test-title", 99.99, "test-note", pq.Array([]string{"food", "snack"}))
	mockedSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"
	mock.ExpectPrepare(regexp.QuoteMeta(mockedSql)).ExpectQuery().WithArgs(1).WillReturnRows(expensesMockRows)

	var latest Expense
	rec := httptest.NewRecorder()

	c := setUpTest(http.MethodGet, uri("expenses", "1"), nil, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := GetExpenseByIdHandler(c)

	json.NewDecoder(rec.Body).Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, 1, latest.ID)
	assert.Equal(t, "test-title", latest.Title)
	assert.Equal(t, 99.99, latest.Amount)
	assert.Equal(t, "test-note", latest.Note)
	assert.Equal(t, []string{"food", "snack"}, latest.Tags)
}

func TestUpdateExpenseByID(t *testing.T) {
	mock := setUpDB(t)

	mockedSql := "UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1"
	mockedRow := sqlmock.NewRows([]string{"id"}).AddRow(1)
	
	mock.ExpectPrepare(regexp.QuoteMeta(mockedSql)).ExpectQuery().WithArgs(1, "Title", 40.5, "notess", pq.Array([]string{"food", "snack"})).WillReturnRows((mockedRow))

	body := bytes.NewBufferString(`{
        "title": "Title",
        "amount": 40.5,
        "note": "notess",
        "tags": ["food", "snack"]
    }`)
	var latest Expense
	rec := httptest.NewRecorder()
	c := setUpTest(http.MethodPut, uri("expenses", "1"), body, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := PutExpenseHandler(c)
	json.NewDecoder(rec.Body).Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, 1, latest.ID)
	assert.Equal(t, "Title", latest.Title)
	assert.Equal(t, 40.5, latest.Amount)
	assert.Equal(t, "notess", latest.Note)
	assert.Equal(t, []string{"food", "snack"}, latest.Tags)
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}
