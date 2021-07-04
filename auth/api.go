package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}

	r.POST("/login", res.login)
}

type resource struct {
	service Service
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r resource) login(c echo.Context) error {
	req := new(LoginDTO)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	t, err := r.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, t)
}
