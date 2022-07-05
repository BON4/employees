package server

import (
	"context"
	"time"

	authHttp "github.com/BON4/employees/internal/authJWT/delivery/http"
	jwtMiddleware "github.com/BON4/employees/internal/authJWT/delivery/http/middleware"
	jwtRepository "github.com/BON4/employees/internal/authJWT/repository"
	empHttp "github.com/BON4/employees/internal/employees/delivery/http"
	treeRepository "github.com/BON4/employees/internal/employees/repository"
	"github.com/BON4/employees/internal/models"
	"github.com/BON4/employees/pkg/jwtService"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var ACCESS_KEY []byte = []byte("123456789")
var REFRESH_KEY []byte = []byte("123456789")
var ACCESS_TIME time.Duration = time.Hour
var REFRESH_TIME time.Duration = time.Hour * 6

func (s *Server) MapHandlers(e *echo.Echo) error {

	e.Use(middleware.Logger())

	eT := models.NewEmpMapTree()
	if err := eT.Load(context.Background(), s.st); err != nil {
		return err
	}

	rep := treeRepository.NewTreeRepo(eT, s.st)

	jwtManagerConf := jwtService.NewJWTConfig(ACCESS_KEY, REFRESH_KEY, ACCESS_TIME, REFRESH_TIME)
	jwtManager := jwtService.NewJWTService(jwtManagerConf)
	jwtRep := jwtRepository.NewJWTRepo(rep, s.st, jwtManager)

	jwtMid := jwtMiddleware.NewJwtMiddleware(jwtManager, s.logger)

	v1 := e.Group("/v1")

	empHandler := empHttp.NewEmployeeHandler(rep, s.logger)
	empGroup := v1.Group("/emp", jwtMid.AuthCheck(), jwtMid.AccessCheck())
	//empGroup := v1.Group("/emp")
	empHttp.NewEmployeeRoutes(empGroup, empHandler)

	authHandler := authHttp.NewAuthHandler(jwtRep, s.logger)
	authGroup := v1.Group("/auth")
	authHttp.NewAuthRoutes(authGroup, authHandler)

	return nil
}
