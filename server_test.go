//go:build integration
// +build integration

package main

import (
	"context"
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

func TestITGetGreeting(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		e.GET("/expenses", TokenCheck(repos.GetExpensesHandler))
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	fmt.Println("-=->", req)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "November 10, 2009")
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	fmt.Println("==> ", req)
	fmt.Println("--> ", resp)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		fmt.Println(byteBody)
		//assert.Equal(t, "Hello, World!", string(byteBody))			//<---
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
