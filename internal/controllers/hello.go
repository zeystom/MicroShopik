package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HelloController struct{}

func NewHelloController() *HelloController {
	return &HelloController{}
}

// SayHello godoc
// @Summary Say hello
// @Description Get a hello message
// @Tags hello
// @Accept json
// @Produce plain
// @Success 200 {string} string "Hello, World!"
// @Router /aboba [get]
func (h *HelloController) SayHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
