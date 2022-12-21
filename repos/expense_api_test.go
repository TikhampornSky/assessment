//go:build integration

package repos

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func seedExpense(t *testing.T) Expense {
	var c Expense
	body := bytes.NewBufferString(`{
		"title": "Sunflower Seed",
		"amount": 11.11,
		"note": "Pink Monday discount10%", 
		"tags": ["food", "snack", "recommennded"]
	}`)

	err := request(http.MethodPost, uri("expenses"), body).Decode((&c))
	if err != nil {
		t.Fatal("can't create expense!: ", err)
	}
	return c
}

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie C++",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	var e Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie C++", e.Title)
	assert.Equal(t, float64(79), e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)
}

func TestGetAllExpenses(t *testing.T) {
	seedExpense(t)

	var expense []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&expense)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(expense), 0)
}

func TestGetExpenseByID(t *testing.T) {
	c := seedExpense(t)

	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.Equal(t, "Sunflower Seed", latest.Title)
	assert.Equal(t, 11.11, latest.Amount)
	assert.Equal(t, "Pink Monday discount10%", latest.Note)
	assert.Equal(t, []string{"food", "snack", "recommennded"}, latest.Tags)
}

func TestUpdateExpenseByID(t *testing.T) {
	c := seedExpense(t)

	body := bytes.NewBufferString(`{
		"title": "Hamtaro Sunflower Seed",
		"amount": 40.5,
		"note": "Pink Monday discount10%", 
		"tags": ["food", "snack", "recommennded"]
	}`)

	var latest Expense
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(c.ID)), body)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.Equal(t, "Hamtaro Sunflower Seed", latest.Title)
	assert.Equal(t, 40.5, latest.Amount)
	assert.Equal(t, "Pink Monday discount10%", latest.Note)
	assert.Equal(t, []string{"food", "snack", "recommennded"}, latest.Tags)
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "November 10, 2009")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
