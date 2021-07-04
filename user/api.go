package user

import (
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}

	r.POST("/user", res.create)
}

type resource struct {
	service Service
}

type UserDTO struct {
	Name     string `json:"name" pg:",notnull" validate:"required"`
	Email    string `json:"email" pg:",notnull" validate:"required,email"`
	Password string `json:"password" pg:",notnull" validate:"required"`
}

func (r resource) create(c echo.Context) error {
	u := new(UserDTO)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	exists, err := r.service.Exists(u.Email)
	if exists {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "This email address is already being used.")
	}
	if err != nil && err != pg.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = r.service.Create(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}
