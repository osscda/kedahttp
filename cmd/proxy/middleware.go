package main

import "github.com/labstack/echo/v4"

func allMiddlewares(counter *reqCounter) echo.MiddlewareFunc {
	countMiddleware := func(fn echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: need to figure out a way to get the increment
			// to happen before fn(w, r) happens below. otherwise,
			// the counter won't get incremented right away and the actual
			// handler will hang longer than it needs to
			go func() {
				counter.inc()
			}()
			defer func() {
				counter.dec()
			}()
			fn(c)
			return nil
		}
	}
	return countMiddleware
}
