package http

import (
	"log"

	"github.com/BON4/employees/internal/employees"
	echo "github.com/labstack/echo/v4"
)

type employeeHandler struct {
	logger *log.Logger
	repo   employees.EmpRepository
}

func (e employeeHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		jsonEmpTree, err := e.repo.Json(c.Request().Context(), uuid)
		if err != nil {
			return echo.ErrNotFound
		}
		e.logger.Printf("%s\n", jsonEmpTree)
		return c.JSONBlob(200, []byte(jsonEmpTree))
	}
}

func NewEmployeeHandler(repo employees.EmpRepository, logger *log.Logger) *employeeHandler {
	return &employeeHandler{
		repo:   repo,
		logger: logger,
	}
}
