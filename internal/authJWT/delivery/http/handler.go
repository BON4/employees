package http

import (
	"encoding/json"
	"log"

	"github.com/BON4/employees/internal/authJWT"
	uerrors "github.com/BON4/employees/internal/errors"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	logger *log.Logger
	repo   authJWT.AuthRepo
}

func (a *authHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		form := LoginForm{}
		err := json.NewDecoder(c.Request().Body).Decode(&form)
		if err != nil {
			a.logger.Println(err.Error())
			return echo.NewHTTPError(500)
		}

		access, refresh, err := a.repo.Login(c.Request().Context(), form.Username, form.Password)
		if err != nil {
			switch err.(type) {
			case uerrors.UserError:
				return echo.NewHTTPError(500, err.Error())
			default:
				a.logger.Println(err.Error())
				return echo.NewHTTPError(500)
			}
		}

		return c.JSON(200, Cookies{
			Access_Token:  access,
			Refresh_Token: refresh,
		})
	}
}

func (a *authHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		form := LoginForm{}
		err := json.NewDecoder(c.Request().Body).Decode(&form)
		if err != nil {
			a.logger.Println(err.Error())
			return echo.NewHTTPError(500)
		}

		access, refresh, err := a.repo.Register(c.Request().Context(), form.Username, form.Password)
		if err != nil {
			switch err.(type) {
			case uerrors.UserError:
				return echo.NewHTTPError(500, err.Error())
			default:
				a.logger.Println(err.Error())
				return echo.NewHTTPError(500)
			}
		}

		return c.JSON(200, Cookies{
			Access_Token:  access,
			Refresh_Token: refresh,
		})
	}
}

func (a *authHandler) Refresh() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("refresh_token")
		if err != nil {
			a.logger.Println(err.Error())
			return echo.NewHTTPError(500)
		}

		access, refresh, err := a.repo.Refresh(c.Request().Context(), cookie.Value)
		if err != nil {
			switch err.(type) {
			case uerrors.UserError:
				return echo.NewHTTPError(500, err.Error())
			default:
				a.logger.Println(err.Error())
				return echo.NewHTTPError(500)
			}
		}

		return c.JSON(200, Cookies{
			Access_Token:  access,
			Refresh_Token: refresh,
		})
	}
}

func NewAuthHandler(repo authJWT.AuthRepo, logger *log.Logger) *authHandler {
	return &authHandler{
		repo:   repo,
		logger: logger,
	}
}
