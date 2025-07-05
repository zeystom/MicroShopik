package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HelloController struct{}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (h *HelloController) SayHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
