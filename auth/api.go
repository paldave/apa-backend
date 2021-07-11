package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type resource struct {
	service Service
}

type loginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	r.POST("/login", res.login)
}

func (r *resource) login(c echo.Context) error {
	req := new(loginDTO)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	t, err := r.service.Authenticate(req)
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, t)
}
