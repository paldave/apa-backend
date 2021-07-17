package auth

import (
	"github.com/labstack/echo/v4"
)

const (
	AuthScheme = "Bearer"
)

type ContextToken struct {
	Id      string
	UserId  string
	Email   string
	IsAdmin bool
}

func Middleware(r Repository, j JWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)

			l := len(AuthScheme)
			if len(auth) > l+1 && auth[:l] == AuthScheme {
				s := auth[l+1:]
				token, claims, err := j.Validate(s)

				if err != nil || !token.Valid {
					return echo.ErrUnauthorized
				}

				exists, err := r.Exists(claims["Id"].(string), claims["UserId"].(string))
				if err != nil || !exists {
					return echo.ErrUnauthorized
				}

				c.Set("authIsAdmin", claims["IsAdmin"])
				c.Set("authUserId", claims["UserId"])

				return next(c)
			}

			return echo.ErrUnauthorized
		}
	}
}
