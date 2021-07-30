package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// TODO: Proper CORS Handling
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Validator = NewValidator(e)
	e.Binder = NewBinder()
	return e
}
