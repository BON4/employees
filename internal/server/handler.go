package server

import (
	empHttp "github.com/BON4/employees/internal/employees/delivery/http"
	"github.com/BON4/employees/internal/employees/repository"
	"github.com/BON4/employees/internal/store"
	echo "github.com/labstack/echo/v4"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	store := store.NewStore()

	rep, err := repository.NewTreeRepoDEBUG(store)
	if err != nil {
		return err
	}

	v1 := e.Group("/v1")

	empHandler := empHttp.NewEmployeeHandler(rep, s.logger)
	empGroup := v1.Group("/emp")
	empHttp.NewEmployeeRoutes(empGroup, *empHandler)

	return nil
}
