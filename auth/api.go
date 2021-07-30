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
	r.POST("/refresh", res.refresh)
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

	c.SetCookie(r.service.BuildCookie("accessToken", t.AccessToken))
	c.SetCookie(r.service.BuildCookie("refreshToken", t.RefreshToken))

	return c.JSON(http.StatusOK, map[string]string{})
}

func (r *resource) logout(c echo.Context) error {
	tokenId := c.Get("authTokenId").(string)
	userId := c.Get("authUserId").(string)
	refreshId := c.Get("authRefreshId").(string)

	if err := r.service.Logout(tokenId, refreshId, userId); err != nil {
		return echo.ErrInternalServerError
	}

	c.SetCookie(&http.Cookie{
		Name:   "accessToken",
		MaxAge: -1,
	})

	c.SetCookie(&http.Cookie{
		Name:   "refreshToken",
		MaxAge: -1,
	})

	return c.JSON(http.StatusOK, map[string]string{})
}

func (r *resource) refresh(c echo.Context) error {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		return echo.ErrUnauthorized
	}

	t, err := r.service.AuthenticateRefresh(cookie.Value)
	if err != nil {
		return echo.ErrUnauthorized
	}

	c.SetCookie(r.service.BuildCookie("accessToken", t.AccessToken))
	c.SetCookie(r.service.BuildCookie("refreshToken", t.RefreshToken))

	return c.JSON(http.StatusOK, t)
}
