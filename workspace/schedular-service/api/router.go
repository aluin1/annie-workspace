package api

import (
	"encoding/json"
	"net/http"

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

func constructIsoError(isoData map[string]string, responseCode string) map[string]string {
	isoData["0"] = "0210"
	isoData["39"] = responseCode
	return isoData
}

// streamJSON stream JSON as response
func StreamJSON(ctx echo.Context, v interface{}) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	ctx.Response().Header().Set("Cache-Control", "no-cache")
	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(v)
}
