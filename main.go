package main

import (
	"os"

	"order-manager/app/infrastructure/http/handler/cookie"
	"order-manager/app/infrastructure/note"
	"order-manager/app/shared/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Handlers Http
	_ = cookie.NewHttpHandler(e, createNoteRepository())

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func createNoteRepository() *note.Repository {
	token := os.Getenv("TOKEN")
	proxy := os.Getenv("HTTP_PROXY")
	dbConfig := &note.DBConfig{
		IDDBOrder: os.Getenv("DB_ID_ORDER"),
	}
	return note.NewListItemRepository(token, dbConfig, utils.GetClient(proxy))

}
