package auth

import (
	"apa-backend/entity"
	"errors"

	jwtgo "github.com/dgrijalva/jwt-go"
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

// Since we attach both access and refresh tokens from browser
// we can reduce extra /refresh endpoint call by refreshing
// tokens straight from the middleware.
func refreshTokens(c echo.Context, s *service) error {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		return errors.New("")
	}

	token, claims, err := s.jwt.Validate(cookie.Value)
	if err != nil || !token.Valid {
		return errors.New("")
	}

	cId := claims["Id"].(string)
	cAID := claims["AccessId"].(string)
	cUID := claims["UserId"].(string)

	exists, err := s.r.Exists(cId, cUID)
	if err != nil || !exists {
		return errors.New("")
	}

	u, err := s.ur.FindById(cUID)
	if err != nil {
		return errors.New("")
	}

	if err := s.r.Delete(cId, cUID); err != nil {
		return errors.New("")
	}

	if err := s.r.Delete(cAID, cUID); err != nil {
		return errors.New("")
	}

	var AuthToken = &entity.AuthToken{}
	tokens, err := s.buildTokens(u, AuthToken)
	if err != nil {
		return errors.New("")
	}

	_, claims, err = s.jwt.Validate(tokens.AccessToken)
	if err != nil {
		return errors.New("")
	}

	c.SetCookie(s.BuildCookie("accessToken", tokens.AccessToken))
	c.SetCookie(s.BuildCookie("refreshToken", tokens.RefreshToken))

	c.Set("authIsAdmin", claims["IsAdmin"])
	c.Set("authUserId", claims["UserId"])
	c.Set("authTokenId", claims["Id"])
	c.Set("authRefreshId", claims["RefreshId"])

	return nil
}

func Middleware(s *service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var isRefreshOnly = false

			cookie, err := c.Cookie("accessToken")
			if err != nil {
				isRefreshOnly = true
				cookie, err = c.Cookie("refreshToken")
				if err != nil {
					return echo.ErrUnauthorized
				}
			}

			token, claims, err := s.jwt.Validate(cookie.Value)

			if isRefreshOnly && err == nil && token.Valid {
				if err := refreshTokens(c, s); err != nil {
					return echo.ErrUnauthorized
				}

				return next(c)
			}

			if ve, ok := err.(*jwtgo.ValidationError); ok {
				if ve.Errors&jwtgo.ValidationErrorExpired != 0 {
					if err := refreshTokens(c, s); err != nil {
						return echo.ErrUnauthorized
					}

					return next(c)
				} else {
					return echo.ErrUnauthorized
				}
			}

			if !token.Valid {
				return echo.ErrUnauthorized
			}

			exists, err := s.r.Exists(claims["Id"].(string), claims["UserId"].(string))
			if err != nil || !exists {
				return echo.ErrUnauthorized
			}

			c.Set("authIsAdmin", claims["IsAdmin"])
			c.Set("authUserId", claims["UserId"])
			c.Set("authTokenId", claims["Id"])
			c.Set("authRefreshId", claims["RefreshId"])

			return next(c)
		}
	}
}

// func Middleware(r Repository, j JWT) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			// auth := c.Request().Header.Get(echo.HeaderAuthorization)

// 			cookie, err := c.Cookie("accessToken")
// 			if err != nil {
// 				return echo.ErrUnauthorized
// 			}

// 			// l := len(AuthScheme)
// 			// if len(auth) > l+1 && auth[:l] == AuthScheme {
// 			// s := auth[l+1:]
// 			fmt.Println(cookie.Value)
// 			token, claims, err := j.Validate(cookie.Value)

// 			if err != nil || !token.Valid {
// 				return echo.ErrUnauthorized
// 			}

// 			exists, err := r.Exists(claims["Id"].(string), claims["UserId"].(string))
// 			if err != nil || !exists {
// 				return echo.ErrUnauthorized
// 			}

// 			c.Set("authIsAdmin", claims["IsAdmin"])
// 			c.Set("authUserId", claims["UserId"])
// 			c.Set("authTokenId", claims["Id"])
// 			c.Set("authRefreshId", claims["RefreshId"])

// 			return next(c)
// 			// }

// 			// return echo.ErrUnauthorized
// 		}
// 	}
// }
