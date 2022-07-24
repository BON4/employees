package http

import (
	"encoding/json"
	"log"

	"github.com/BON4/employees/internal/employees"
	uerrors "github.com/BON4/employees/internal/errors"
	echo "github.com/labstack/echo/v4"
)

type employeeHandler struct {
	logger *log.Logger
	repo   employees.EmpRepository
}

func (e *employeeHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		jsonEmpTree, err := e.repo.Json(c.Request().Context(), uuid)
		if err != nil {
			switch err.(type) {
			case uerrors.UserError:
				return echo.NewHTTPError(500, err.Error())
			default:
				e.logger.Println(err.Error())
				return echo.NewHTTPError(500)
			}
		}
		return c.JSONBlob(200, []byte(jsonEmpTree))
	}
}

func (e *employeeHandler) Move() echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")

		josnForm := MoveForm{}

		if err := json.NewDecoder(c.Request().Body).Decode(&josnForm); err != nil {
			e.logger.Println(err.Error())
			return echo.NewHTTPError(500)
		}

		if err := e.repo.Move(c.Request().Context(), uuid, josnForm.FromUUID, josnForm.ToUUID); err != nil {
			switch err.(type) {
			case uerrors.UserError:
				return echo.NewHTTPError(500, err.Error())
			default:
				e.logger.Println(err.Error())
				return echo.NewHTTPError(500)
			}
		}

		return c.JSON(200, []byte(""))
	}
}

func NewEmployeeHandler(repo employees.EmpRepository, logger *log.Logger) *employeeHandler {
	return &employeeHandler{
		repo:   repo,
		logger: logger,
	}
}
