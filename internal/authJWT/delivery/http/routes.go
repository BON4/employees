package http

import "github.com/labstack/echo/v4"

func NewAuthRoutes(group *echo.Group, h *authHandler) {
	group.POST("/login", h.Login())
	group.POST("/register", h.Register())
	group.POST("/refresh", h.Refresh())
}
