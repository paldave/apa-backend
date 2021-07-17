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

func RegisterHandlers(r *echo.Group, middleware echo.MiddlewareFunc, service Service) {
	res := &resource{service}

	r.POST("/login", res.login)
	r.POST("/logout", res.logout, middleware)
}

func (r *resource) login(c echo.Context) error {
	req := new(loginDTO)
	if err := c.Bind(req); err != nil {
		return echo.ErrInternalServerError
	}

	t, err := r.service.Authenticate(req)
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, t)
}

func (r *resource) logout(c echo.Context) error {
	tokenId := c.Get("authTokenId").(string)
	userId := c.Get("authUserId").(string)

	if err := r.service.Logout(tokenId, userId); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{})
}
