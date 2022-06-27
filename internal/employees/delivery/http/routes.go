package http

import "github.com/labstack/echo/v4"

func NewEmployeeRoutes(group *echo.Group, h employeeHandler) {
	group.GET("/:uuid/list", h.List())
}
