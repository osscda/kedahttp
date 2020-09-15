package main

import (
	"github.com/labstack/echo/v4"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Request().Header.Set("X-Echo-Error", "true")
	c.Logger().Error(err)
}
