package cookie

import (
	"context"
	"fmt"
	"net/http"

	"order-manager/app/domain/repository"

	"github.com/labstack/echo/v4"
)

type HttpHandler struct {
	cookieRepository repository.Cookie
}

func NewHttpHandler(e *echo.Echo, cookieRepository repository.Cookie) *HttpHandler {
	h := &HttpHandler{cookieRepository: cookieRepository}
	// Routes
	e.GET("/cookie/add", h.add)
	e.GET("/cookie/sell-one", h.sellOne)
	//TODO: implement to future
	//e.GET("/cookie/sell/{qua}", h.sellOne)
	e.GET("/cookie", h.current)

	return h
}

type Response struct {
	Message       string `json:"message"`
	CurrentCookie int    `json:"current_cookie"`
}

func (h *HttpHandler) add(c echo.Context) error {
	u := &Response{
		Message:       "Tu galleta fue añadida, Esta data es fake",
		CurrentCookie: 14,
	}
	return c.JSON(http.StatusOK, u)
}

func (h *HttpHandler) sellOne(c echo.Context) error {
	currentCookie, err := h.cookieRepository.SellOne(context.Background())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Message:       fmt.Sprintf("Error al restar galleta: %v", err),
			CurrentCookie: 0,
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Message:       "Tu galleta fue restada",
		CurrentCookie: currentCookie,
	})
}

func (h *HttpHandler) current(c echo.Context) error {
	currentCookie, err := h.cookieRepository.Current(context.Background())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Message:       fmt.Sprintf("Error al obtener las galletas actuales: %v", err),
			CurrentCookie: 0,
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Message:       "Galletas actuales",
		CurrentCookie: currentCookie,
	})
}
