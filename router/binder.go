package router

import (
	"github.com/labstack/echo/v4"
)

func NewBinder() *CustomBinder {
	return &CustomBinder{b: &echo.DefaultBinder{}}
}

type CustomBinder struct {
	b echo.Binder
}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) error {
	if err := cb.b.Bind(i, c); err != nil && err != echo.ErrUnsupportedMediaType {
		return err
	}

	return c.Validate(i)
}
