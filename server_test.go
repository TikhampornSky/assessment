//go:build integration
// +build integration

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/TikhampornSky/assessment/repos"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const serverPort = 2565

var expense_global repos.Expense

// To run these tests in this file use command docker compose up --build go_test

func startServer() *echo.Echo {
	eh := echo.New()
	go func(e *echo.Echo) {
		e.GET("/expenses", TokenCheck(repos.GetExpensesHandler))
		e.POST("/expenses", TokenCheck(repos.CreateExpenseHandler))
		e.GET("/expenses/:id", TokenCheck(repos.GetExpenseHandler))
		e.PUT("/expenses/:id", TokenCheck(repos.PutExpenseHandler))
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("go:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println("(", err, ")", "Please wait a minute..")
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	return eh
}

func closeServer(eh *echo.Echo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := eh.Shutdown(ctx)
	return err
}

func request(req *http.Request) *http.Response {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "November 10, 2009")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	return resp
}

func TestGetEveryExpenses(t *testing.T) {

	eh := startServer()

	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://go:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)

	resp := request(req)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Greater(t, len(string(byteBody)), 0)
	}

	err = closeServer(eh)
	assert.NoError(t, err)
}

func TestPostExpenses(t *testing.T) {

	eh := startServer()

	reqBody := bytes.NewBufferString(`{
		"title": "strawberry smoothie C++",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://go:%d/expenses", serverPort), reqBody)
	assert.NoError(t, err)

	resp := request(req)

	var expense repos.Expense
	json.NewDecoder(resp.Body).Decode(&expense)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEqual(t, 0, expense.ID)
		assert.Equal(t, "strawberry smoothie C++", expense.Title)
		assert.Equal(t, float64(79), expense.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", expense.Note)
		assert.Equal(t, []string{"food", "beverage"}, expense.Tags)
	}

	err = closeServer(eh)
	assert.NoError(t, err)

	expense_global = expense
}

func TestGetExpense(t *testing.T) {

	eh := startServer()

	TestPostExpenses(t)
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://go:%d/expenses/%d", serverPort, expense_global.ID), strings.NewReader(reqBody))
	assert.NoError(t, err)

	resp := request(req)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Greater(t, len(string(byteBody)), 0)
		assert.Equal(t, "strawberry smoothie C++", expense_global.Title)
		assert.Equal(t, float64(79), expense_global.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", expense_global.Note)
		assert.Equal(t, []string{"food", "beverage"}, expense_global.Tags)
	}

	err = closeServer(eh)
	assert.NoError(t, err)
}

func TestPutExpense(t *testing.T) {

	eh := startServer()

	TestPostExpenses(t)
	reqBody := bytes.NewBufferString(`{
		"title": "strawberry smoothie From Taiwan",
		"amount": 99,
		"note": "night market promotion",
		"tags": ["food", "beverage", "new release"]
	}`)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://go:%d/expenses/%d", serverPort, expense_global.ID), reqBody)
	assert.NoError(t, err)
	
	resp := request(req)

	var expense repos.Expense
	json.NewDecoder(resp.Body).Decode(&expense)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "strawberry smoothie From Taiwan", expense.Title)
		assert.Equal(t, float64(99), expense.Amount)
		assert.Equal(t, "night market promotion", expense.Note)
		assert.Equal(t, []string{"food", "beverage", "new release"}, expense.Tags)
	}

	err = closeServer(eh)
	assert.NoError(t, err)
}