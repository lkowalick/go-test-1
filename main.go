package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lkowalick/go-test-1/cloudsql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dbPool, err := cloudsql.ConnectUnixSocket()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	e.GET("/", func(c echo.Context) error {
		var greeting string
		err = dbPool.QueryRow("select 'Hello, world!'").Scan(&greeting)
		if err != nil {
			e.Logger.Fatal(err)
			return c.HTML(http.StatusInternalServerError, "Hello, Docker! (NO DB!) ")
		}
		return c.HTML(http.StatusOK, fmt.Sprintf("Hello, Docker! %s\n", greeting))
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// Simple implementation of an integer minimum
// Adapted from: https://gobyexample.com/testing-and-benchmarking
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
