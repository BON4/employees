package middleware

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/BON4/employees/internal/models"
	"github.com/BON4/employees/pkg/jwtService"
	"github.com/labstack/echo/v4"
)

type jwtMiddleware struct {
	logger     *log.Logger
	jwtManager *jwtService.JWTService
}

func (j *jwtMiddleware) AuthCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			j.logger.Println("Got AuthCheck")
			cookie, err := c.Cookie("access_token")
			if err != nil {
				j.logger.Println(err.Error())
				return err
			}

			tokData, err := j.jwtManager.VerifyAcess(cookie.Value)
			if err != nil {
				j.logger.Println(err.Error())
				return err
			}

			emp := &models.Employee{}
			if empByte, ok := tokData["Emp"]; ok {
				if empJSON, err := json.Marshal(empByte); err == nil {
					if err := json.Unmarshal(empJSON, emp); err != nil {
						j.logger.Println(err.Error())
						return err
					}
				}
			} else {
				return errors.New("invalid token data")
			}

			c.Set("Emp", emp)

			return next(c)
		}
	}
}

func (j *jwtMiddleware) AccessCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			j.logger.Println("Got AccessCheck")
			if emp, ok := c.Get("Emp").(*models.Employee); ok {
				if emp.UUID == c.Param("uuid") {
					return next(c)
				}
			}

			return errors.New("access denied")
		}
	}
}

func NewJwtMiddleware(jwtManager *jwtService.JWTService, logger *log.Logger) *jwtMiddleware {
	return &jwtMiddleware{
		jwtManager: jwtManager,
		logger:     logger,
	}
}
