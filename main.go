package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/cookie/add", cookieAdd)
	e.GET("/cookie/subtract", cookieSubtract)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

type Response struct {
	Message       string `json:"message"`
	CurrentCookie int    `json:"current_cookie"`
}

func cookieAdd(c echo.Context) error {
	u := &Response{
		Message:       "Tu galleta fue añadida",
		CurrentCookie: 14,
	}
	return c.JSON(http.StatusOK, u)
}

func cookieSubtract(c echo.Context) error {
	u := &Response{
		Message:       "Tu galleta fue restada",
		CurrentCookie: 2,
	}
	return c.JSON(http.StatusOK, u)
}
