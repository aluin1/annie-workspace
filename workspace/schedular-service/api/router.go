package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/resty.v1"
)

var caseRest *resty.Client

// ConfigRouter configure API router
func ConfigRouter() *echo.Echo {

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} ${uri} | status=${status} | ${latency_human}\n",
	}))

	e.GET("/ping", ping)

	return e
}

func ping(ctx echo.Context) (err error) {
	return nil
}
