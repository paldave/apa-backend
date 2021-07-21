package user

import (
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
)

type resource struct {
	service Service
}

type userDTO struct {
	Name     string `json:"name" pg:",notnull" validate:"required"`
	Email    string `json:"email" pg:",notnull" validate:"required,email"`
	Password string `json:"password" pg:",notnull" validate:"required"`
}

func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	r.POST("/user", res.create)
}

func (r resource) create(c echo.Context) error {
	u := new(userDTO)
	if err := c.Bind(u); err != nil {
		return echo.ErrInternalServerError
	}

	exists, err := r.service.Exists(u.Email)
	if exists {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "This email address is already being used.")
	}
	if err != nil && err != pg.ErrNoRows {
		return echo.ErrInternalServerError
	}

	_, err = r.service.Create(u)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, u)
}
